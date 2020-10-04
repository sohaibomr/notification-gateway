package models

// NotificationRequest ...
type NotificationRequest struct {
	Type      string   `json:"type" binding:"required,oneof=group user"` //group or custom/personalized
	SendVia   string   `json:"sendVia" binding:"required,oneof=sms push"`
	Message   string   `json:"message" binding:"required"`
	Category  string   `json:"category" binding:"required"` // promo code, destination reached, captain waiting etc
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
}
