package app

import (
	"sender/app/command"
	"sender/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SentBroadcast               command.SendNotificationBroadcastHandler
	SentNotification            command.SendNotificationHandler
	LogInAndGetOldNotifications command.LogInAndGetOldNotificationsHandler
	UpdateNotificationStatus    command.UpdateNotificationStatusHandler
}

type Queries struct {
	ClientList       query.ClientListHandler
	NotificationList query.NotificationListHandler
}
