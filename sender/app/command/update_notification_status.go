package command

import (
	"context"
	"log"
	"sender/app/decorator"
	"sender/domain/client"
	"sender/domain/notification"

	"github.com/google/uuid"
)

type UpdateNotificationStatusHandler decorator.CommandHandler[UpdateNotification]

type updateNotificationStatusHandler struct {
	server socketSender
	nRepo  notification.Repository
	cRepo  client.Repository
}

type UpdateNotification struct {
	Id         uuid.UUID
	ClientName string
	Status     notification.NotificationStatus
}

func (h updateNotificationStatusHandler) Handle(ctx context.Context, command UpdateNotification) error {
	clt, err := h.cRepo.GetClientByName(ctx, command.ClientName)
	if err != nil {
		log.Println(err)
		return err
	}
	err = h.nRepo.UpdateNotification(ctx, command.Id, clt.ID(),
		func(n *notification.Notification) (*notification.Notification, error) {
			n.UpdateStatus(command.Status)
			return n, nil
		},
	)
	return err
}
func NewUpdateNotificationStatusHandler(
	server socketSender,
	cRepo client.Repository,
	nRepo notification.Repository) UpdateNotificationStatusHandler {
	if server == nil {
		panic("empty socket server")
	}
	if cRepo == nil {
		panic("empty client repository")
	}
	if nRepo == nil {
		panic("empty notification repository")
	}
	return decorator.NewCommandHadlerWithDefaultDecorators[UpdateNotification](updateNotificationStatusHandler{
		server: server,
		cRepo:  cRepo,
		nRepo:  nRepo,
	})
}
