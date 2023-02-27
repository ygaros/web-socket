package command

import (
	"context"
	"log"
	"sender/app/decorator"
	"sender/domain"
	"sender/domain/client"
	"sender/domain/notification"
)

type LogInAndGetOldNotificationsHandler decorator.CommandHandler[string]

type logInAndGetOldNotificationsHandler struct {
	server socketSender
	cRepo  client.Repository
	nRepo  notification.Repository
}

func (h logInAndGetOldNotificationsHandler) Handle(ctx context.Context, username string) error {
	clt, err := h.cRepo.GetClientByName(ctx, username)
	if err != nil {
		c, err := client.NewClient(username)
		if err != nil {
			log.Println(err)
			return err
		}
		clt, err = h.cRepo.SaveClient(ctx, c)
		if err != nil {
			log.Println(err)
		}
	}

	notifications, err := h.nRepo.GetAllNotificationByClientId(ctx, clt.ID())
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("notifications length %d\n", len(notifications))
	var parsedToSend []domain.NotificationToSend
	for _, n := range notifications {
		// err := h.server.SendToRecipient(domain.NotificationToSend{
		// 	Id:      n.ID(),
		// 	Message: n.Message(),
		// 	Status:  n.Status().String(),
		// }, username)
		if err != nil {
			log.Println(err)
			continue
		}
		parsedToSend = append(parsedToSend, domain.NotificationToSend{
			Id:      n.ID(),
			Message: n.Message(),
			Status:  n.Status().String(),
		})
	}
	h.server.SendMultipleToRecipient(parsedToSend, username)
	return err

}

func NewLogInAndGetOldNotificationsHandler(
	server socketSender,
	cRepo client.Repository,
	nRepo notification.Repository) LogInAndGetOldNotificationsHandler {
	if server == nil {
		panic("empty socket server")
	}
	if cRepo == nil {
		panic("empty client repository")
	}
	if nRepo == nil {
		panic("empty notification repository")
	}
	return decorator.NewCommandHadlerWithDefaultDecorators[string](logInAndGetOldNotificationsHandler{
		server: server,
		cRepo:  cRepo,
		nRepo:  nRepo,
	})
}
