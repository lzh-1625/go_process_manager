package model



type WsMessage struct {
	MessageType string `json:"messageType"`
	Content     string `json:"content"`
}
