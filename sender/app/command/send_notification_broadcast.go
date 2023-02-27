package command

import (
	"context"
	"sender/app/decorator"
	"sender/domain"
)

type SendNotificationBroadcastHandler decorator.CommandHandler[SendNotificationBroadcast]

type sendNotificationBroadcastHandler struct {
	server socketSender
}

type SendNotificationBroadcast struct {
	Message string
	Sender  string
}

func (h sendNotificationBroadcastHandler) Handle(ctx context.Context, command SendNotificationBroadcast) error {

	return h.server.SendToAllClients(domain.NotificationToSendBroadcast{
		Message: command.Message,
		Sender:  command.Sender,
	})
}

func NewSendNotificationBroadcastHandler(server socketSender) SendNotificationBroadcastHandler {
	if server == nil {
		panic("empty socket server")
	}
	return decorator.NewCommandHadlerWithDefaultDecorators[SendNotificationBroadcast](sendNotificationBroadcastHandler{
		server: server,
	})
}
