package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/andrewawni/notification_service/common/integration"
	"github.com/andrewawni/notification_service/common/messaging"
	"github.com/andrewawni/notification_service/common/models"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

var client messaging.IMessagingClient
var rdb *redis.Client
var ctx context.Context

const processedNotificationsQueueName = "notification_service:dispatch_notifications_jobs"
const smsBucketKey = "notification_service:sms_bucket"
const emailBucketKey = "notification_service:email_bucket"
const pushNotificationBucketKey = "notification_service:push_notifications_bucket"

var bucketsLimitsPerMinute = map[string]int{
	smsBucketKey:              60,
	emailBucketKey:            40,
	pushNotificationBucketKey: 20,
}

func processedNotificationsWorker(d amqp.Delivery) {
	// get user by id, select fields in tags
	notification := models.ProcessedNotification{}
	err := json.Unmarshal(d.Body, &notification)
	if err != nil {
		log.Printf("Failed to process message")
		return
	}

	switch notification.Method {
	case "sms":
		if isAllowed(smsBucketKey) {
			integration.SendSMS(notification)
		} else {
			client.PublishOnQueue([]byte(d.Body), processedNotificationsQueueName)
		}
	case "email":
		if isAllowed(emailBucketKey) {
			integration.SendEmail(notification)
		} else {
			client.PublishOnQueue([]byte(d.Body), processedNotificationsQueueName)
		}
	case "push_notification":
		if isAllowed(pushNotificationBucketKey) {
			integration.SendPushNotification(notification)
		} else {
			client.PublishOnQueue([]byte(d.Body), processedNotificationsQueueName)
		}
	}
}

func isAllowed(key string) bool {

	pipe := rdb.Pipeline()
	decr := pipe.Decr(ctx, key)
	pipe.Exec(ctx)
	log.Println(key, decr.Val())
	return decr.Val() > 0
}

func refillBuckets() {
	for key, val := range bucketsLimitsPerMinute {
		log.Println("refilling ", key, "with ", val)
		err := rdb.Set(ctx, key, val, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	if rdb == nil {
		panic("Can't connect to redis cluster")
	}
	ctx = context.Background()

	refillBuckets()
	ticker := time.NewTicker(1 * time.Minute)

	forever := make(chan bool)
	go func() {
		for {
			select {
			case <-forever:
				return
			case t := <-ticker.C:
				log.Println("Tick at", t)
				refillBuckets()
			}
		}
	}()

	client = &messaging.MessagingClient{}
	client.ConnectToBroker(os.Getenv("RABBITMQ_URL"))
	err := client.SubscribeToQueue(processedNotificationsQueueName, "processed_notifications_worker", 5, processedNotificationsWorker)
	if err != nil {
		panic(err.Error())
	}

	<-forever
}
