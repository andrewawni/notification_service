package models

type User struct {
	Id           string `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	DeviceToken  string `json:"device_token"`
	Locale       string `json:"locale"`
}
