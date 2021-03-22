package models

type Notification struct {
	ID      string `json:"notification_id"`
	Content string `json:"content"`
	Method  string `json:"method"`
}

type SingleNotification struct {
	Notification
	UserID              string   `json:"user_id"`
	PersonalizationTags []string `json:"personalization_tags"`
}

type GroupNotification struct {
	Notification
	GroupID string `json:"group_id"`
}

type ProcessedNotification struct {
	Notification
	Targets []string `json:"targets"`
}
