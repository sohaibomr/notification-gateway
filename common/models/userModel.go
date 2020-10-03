package models

type UserInfo struct {
	ID                    string `json:"userId"`
	Name                  string `json:"name"`
	MobileNo              string `json:"mobileNo"`
	PushNotificationToken string `json:"pushNotificationToken"`
	Email                 string `json:"email"`
}

var UsersMap = map[string]UserInfo{
	"1": UserInfo{ID: "1", Name: "Kashif", MobileNo: "123456789", PushNotificationToken: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Email: "kashif@xyz.com"},
	"2": UserInfo{ID: "2", Name: "John", MobileNo: "123456789", PushNotificationToken: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Email: "John@xyz.com"},
	"3": UserInfo{ID: "3", Name: "Hilary", MobileNo: "123456789", PushNotificationToken: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Email: "Hilary@xyz.com"},
	"4": UserInfo{ID: "4", Name: "Tom", MobileNo: "123456789", PushNotificationToken: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Email: "Tom@xyz.com"},
}

type UserMsg struct {
	Message    string   `json:"message"`
	UserDetail UserInfo `json:"userDetail"`
	SendVia    string   `json:"sendVia"`
}

type Groups struct {
	GroupID string     `json:"groupId"`
	Name    string     `json:"name"`
	Users   []UserInfo `json:"users"`
}

var GroupUsers = map[string][]UserInfo{
	"1": []UserInfo{UsersMap["1"], UsersMap["2"], UsersMap["3"]},
	"2": []UserInfo{UsersMap["1"], UsersMap["2"], UsersMap["3"]},
	"3": []UserInfo{UsersMap["1"], UsersMap["2"], UsersMap["3"]},
	"4": []UserInfo{UsersMap["4"], UsersMap["2"], UsersMap["3"]},
}
