package entity

type MessageSend struct {
	To      string `json:"to" validate:"required"`
	Message string `json:"message" validate:"required"`
	ImageID string `json:"image_id"`
}
