package model

// NotificationRequest model for notification request
type NotificationRequest struct {
	Token string `json:"token"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// RegisterTokenRequest model for register firebase token
type RegisterTokenRequest struct {
	FirebaseToken string `form:"firebase_token" json:"firebase_token"`
}
