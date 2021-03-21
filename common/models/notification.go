package models

type notification struct {
	ID      string `json:"notification_id"`
	Content string `json:"content"`
	Method  string `json:"method"`
}

type SingleNotification struct {
	notification
	UserID              string   `json:"user_id"`
	PersonalizationTags []string `json:"personalization_tags"`
}

type GroupNotification struct {
	notification
	GroupID string `json:"group_id"`
}

type GroupNotificationDelivery struct {
	GroupNotification
	BatchID  string   `json:"batch_id"`
	UsersIDs []string `json:"users_ids"`
}
