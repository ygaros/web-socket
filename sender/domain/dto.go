package domain

import (
	"sender/domain/client"
	"sender/domain/notification"
	"time"

	"github.com/google/uuid"
)

type NotificationRequest struct {
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
	Sender     string   `json:"sender"`
}
type NotificationToSend struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
}
type NotificationToSendBroadcast struct {
	Message string
	Sender  string
}
type NotificationToChangeStatus struct {
	Id         uuid.UUID `json:"id"`
	ClientName string    `json:"clientName"`
	Status     string    `json:"status"`
}
type ClientDTO struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func ToClientDTO(c client.Client) ClientDTO {
	return ClientDTO{
		Id:   c.ID(),
		Name: c.Name(),
	}
}

type NotificationDTO struct {
	Id uuid.UUID `json:"id"`

	Message string `json:"message"`
	Status  string `json:"status"`

	// Sender      string `json:"sender"`
	RecipientId uuid.UUID `json:"recipient_id"`

	CreationDate time.Time `json:"creation_date"`
}

func ToNotificationDTO(n notification.Notification) NotificationDTO {
	return NotificationDTO{
		Id:           n.ID(),
		Message:      n.Message(),
		Status:       n.Status().String(),
		RecipientId:  n.RecipientId(),
		CreationDate: n.CreationDate(),
	}
}
