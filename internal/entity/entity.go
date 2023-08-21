package entity

type StreamOnlineNotification struct {
	Subscription Subscription `json:"subscription"`
	Event        Event        `json:"event"`
}

type Subscription struct {
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	Cost      int       `json:"cost"`
	CreatedAt string    `json:"created_at"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
}

type Condition struct {
	BroadcasterUserId string `json:"broadcaster_user_id"`
}

type Event struct {
	Id                   string `json:"id"`
	BroadcasterUserId    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Type                 string `json:"type"`
	StartedAt            string `json:"started_at"`
}


type UserInfo struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

type UserCollection struct {
	Data []UserInfo `json:"data"`
}