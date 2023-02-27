package notification

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetAllNotifications(ctx context.Context) ([]Notification, error)
	GetNotification(ctx context.Context, notificationId uuid.UUID, clientId uuid.UUID) (*Notification, error)
	GetAllNotificationByClientId(ctx context.Context, clientId uuid.UUID) ([]Notification, error)
	UpdateNotification(
		ctx context.Context,
		notificationId uuid.UUID,
		clientId uuid.UUID,
		updateFn func(n *Notification) (*Notification, error),
	) error
	SaveNotification(ctx context.Context, notification *Notification) (*Notification, error)
}
