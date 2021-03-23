package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andrewawni/notification_service/common/integration"
	"github.com/andrewawni/notification_service/common/messaging"
	"github.com/andrewawni/notification_service/common/models"
	"github.com/streadway/amqp"
)

var client messaging.IMessagingClient

const singleNotificationQueueName = "notification_service:assemble_single_notifications_jobs"
const groupNotificationQueueName = "notification_service:assemble_group_notifications_jobs"
const processedNotificationsQueueName = "notification_service:dispatch_notifications_jobs"
const groupNotificationsBatchSize = 4

var methodToTarget = map[string]string{
	"sms":               "mobile_number",
	"email":             "email",
	"push_notification": "device_token",
}

func singleNotificationsWorker(d amqp.Delivery) {
	// get user by id, select fields in tags
	notification := models.SingleNotification{}
	err := json.Unmarshal(d.Body, &notification)
	if err != nil {
		log.Printf("Failed to process message")
		return
	}

	fields := []string{"locale"}
	fields = append(fields, methodToTarget[notification.Method])
	fields = append(fields, notification.PersonalizationTags...)

	attributes := integration.GetUserAttributes(notification.UserID, fields)
	// format notification content
	content := notification.Content
	for _, attr := range notification.PersonalizationTags {
		pattern := fmt.Sprint("%", attr, "%")
		content = strings.Replace(content, pattern, attributes[attr], -1)
	}
	notification.Content = integration.TranslateContentByLocale(content, attributes["locale"])
	processedNotification := models.ProcessedNotification{Notification: notification.Notification}
	processedNotification.Targets = []string{attributes[methodToTarget[notification.Method]]}
	// push on queue
	payload, _ := json.Marshal(&processedNotification)
	client.PublishOnQueue([]byte(payload), processedNotificationsQueueName)
}

func groupNotificationsWorker(d amqp.Delivery) {
	notification := models.GroupNotification{}
	err := json.Unmarshal(d.Body, &notification)
	if err != nil {
		log.Printf("Failed to process message")
		return
	}
	usersIDs := integration.GetUsersIDsByGroupID(notification.GroupID)
	batchSize := groupNotificationsBatchSize
	targetField := methodToTarget[notification.Method]
	for i := 0; i < len(usersIDs); i += batchSize {
		end := i + batchSize
		if end > len(usersIDs) {
			end = len(usersIDs)
		}
		batchUsersIDs := usersIDs[i:end]
		processedNotification := models.ProcessedNotification{Notification: notification.Notification}
		targets := []string{}
		for _, userID := range batchUsersIDs {
			attr := integration.GetUserAttributes(userID, []string{targetField})
			targets = append(targets, attr[targetField])
		}
		processedNotification.Targets = targets
		payload, _ := json.Marshal(&processedNotification)
		client.PublishOnQueue([]byte(payload), processedNotificationsQueueName)
	}
}

func main() {
	client = &messaging.MessagingClient{}
	client.ConnectToBroker(os.Getenv("RABBITMQ_URL"))
	forever := make(chan bool)
	err := client.SubscribeToQueue(singleNotificationQueueName, "single_notifications_worker", 1, singleNotificationsWorker)
	if err != nil {
		panic(err.Error())
	}
	err = client.SubscribeToQueue(groupNotificationQueueName, "group_notifications_worker", 1, groupNotificationsWorker)
	if err != nil {
		panic(err.Error())
	}
	<-forever
}
