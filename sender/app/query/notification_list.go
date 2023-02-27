package query

import (
	"context"
	"sender/app/decorator"
	"sender/domain"
	"sender/domain/notification"
)

type NotificationListHandler decorator.QueryHandler[interface{}, []domain.NotificationDTO]
type notificationListHandler struct {
	nRepo notification.Repository
}

func (h notificationListHandler) Handle(ctx context.Context, _ interface{}) ([]domain.NotificationDTO, error) {
	notifications, err := h.nRepo.GetAllNotifications(ctx)
	if err != nil {
		return nil, err
	}
	var parsedNotifications []domain.NotificationDTO
	for _, n := range notifications {
		dto := domain.ToNotificationDTO(n)
		parsedNotifications = append(parsedNotifications, dto)
	}

	return parsedNotifications, nil
}
func NewNotificationListHandler(nRepo notification.Repository) NotificationListHandler {
	if nRepo == nil {
		panic("empty notification repository")
	}
	return decorator.NewQueryHandlerWithDefaultDecorators[interface{}, []domain.NotificationDTO](notificationListHandler{nRepo})
}
