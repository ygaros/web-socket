package command

import (
	"context"
	"log"
	"sender/app/decorator"
	"sender/domain"
	"sender/domain/client"
	"sender/domain/notification"
)

type SendNotificationHandler decorator.CommandHandler[SendNotification]

type socketSender interface {
	SendToRecipient(message domain.NotificationToSend, recipient string) error
	SendMultipleToRecipient(messages []domain.NotificationToSend, recipient string) error
	SendToAllClients(message domain.NotificationToSendBroadcast) error
}

type sendNotificationHandler struct {
	server socketSender
	cRepo  client.Repository
	nRepo  notification.Repository
}
type SendNotification struct {
	Message    string
	Sender     string
	Recipients []string
}

func (h sendNotificationHandler) Handle(ctx context.Context, command SendNotification) (err error) {
	//TODO rethink about notification and its id
	// and decide what to do with the continues in this for loop
	var notfs []*notification.Notification
	for _, recipient := range command.Recipients {
		client, err := h.cRepo.GetClientByName(ctx, recipient)
		if err != nil {
			log.Println(err)
			continue
		}
		n, err := notification.NewNotification(command.Message, command.Sender, client.ID())
		if err != nil {
			log.Println(err)
			continue
		}

		notf, err := h.nRepo.SaveNotification(ctx, n)
		if err != nil {
			log.Println(err)
			continue
		}
		notfs = append(notfs, notf)
	}
	for _, n := range notfs {
		client, err := h.cRepo.GetClient(ctx, n.RecipientId())
		if err != nil {
			log.Println(err)
			continue
		}
		err = h.server.SendToRecipient(domain.NotificationToSend{
			Id:      n.ID(),
			Message: n.Message(),
			Status:  n.Status().String(),
		}, client.Name())
	}

	return err
}

func NewSendNotificationHandler(server socketSender, cRepo client.Repository, nRepo notification.Repository) SendNotificationHandler {
	if server == nil {
		panic("empty socket server")
	}
	if cRepo == nil {
		panic("empty client repository")
	}
	if nRepo == nil {
		panic("empty notification repository")
	}
	return decorator.NewCommandHadlerWithDefaultDecorators[SendNotification](sendNotificationHandler{
		server: server,
		cRepo:  cRepo,
		nRepo:  nRepo,
	})
}
