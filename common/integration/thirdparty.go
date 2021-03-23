package integration

import (
	"log"
	"time"

	"github.com/andrewawni/notification_service/common/models"
)

func SendSMS(notification models.ProcessedNotification) {
	time.Sleep(2 * time.Second)
	log.Print("[SMS] sent ", notification.Content, "To ", notification.Targets)
}

func SendPushNotification(notification models.ProcessedNotification) {
	time.Sleep(2 * time.Second)
	log.Print("[PUSH NOTIFICATION] sent ", notification.Content, "To ", notification.Targets)
}

func SendEmail(notification models.ProcessedNotification) {
	time.Sleep(2 * time.Second)
	log.Print("[EMAIL] sent ", notification.Content, "To ", notification.Targets)
}
