package integration

import "github.com/andrewawni/notification_service/common/models"

var users = map[string]models.User{
	"1": {
		Id:           "1",
		Name:         "Andrew",
		Email:        "andrew@gmail.com",
		MobileNumber: "123123124",
		DeviceToken:  "xxxxxxxxx",
		Locale:       "en-US",
	},
	"2": {
		Id:           "2",
		Name:         "Mina",
		Email:        "mina@gmail.com",
		MobileNumber: "123123124",
		DeviceToken:  "xxxxxxxxx",
		Locale:       "en-UK",
	},
	"3": {
		Id:           "3",
		Name:         "Mustafa",
		Email:        "mustafa@gmail.com",
		MobileNumber: "123123124",
		DeviceToken:  "xxxxxxxxx",
		Locale:       "ar-EG",
	},
	"4": {
		Id:           "4",
		Name:         "Sara",
		Email:        "sara@gmail.com",
		MobileNumber: "123123124",
		DeviceToken:  "xxxxxxxxx",
		Locale:       "ar-EG",
	},
	"5": {
		Id:           "5",
		Name:         "Sandy",
		Email:        "sandy@gmail.com",
		MobileNumber: "123123124",
		DeviceToken:  "xxxxxxxxx",
		Locale:       "de-DE",
	},
}

var groups = map[string][]string{
	"1": {"1", "2", "3"},
	"2": {"1", "2", "3", "4", "5"},
	"3": {"3", "5"},
}

func GetUserAttributes(userID string, attributes []string) map[string]string {
	user := users[userID]
	attributesMap := map[string]bool{}
	attributesValues := map[string]string{}

	for _, val := range attributes {
		attributesMap[val] = true
	}

	// TODO convert user object into hash and use dynamic attributes
	if attributesMap["name"] {
		attributesValues["name"] = user.Name
	}
	if attributesMap["email"] {
		attributesValues["email"] = user.Email
	}
	if attributesMap["mobile_number"] {
		attributesValues["mobile_number"] = user.MobileNumber
	}
	if attributesMap["device_token"] {
		attributesValues["device_token"] = user.DeviceToken
	}
	if attributesMap["Locale"] {
		attributesValues["locale"] = user.Locale
	}

	return attributesValues
}

func GetUsersIDsByGroupID(groupID string) []string {
	return groups[groupID]
}
