package main

import (
	"database/sql"
	"log"

	"sender/adapters"
	"sender/app"
	"sender/app/command"
	"sender/app/query"
	"sender/domain/notification"
	"sender/ports"
)

func main() {
	host := ""
	server := ports.NewSocketServer(host, 9876)
	app := NewApplication(server)
	server.SetApp(&app)
	go server.Serve()
	httpServer := ports.NewRestServer(host, 8080, &app)
	log.Fatalln(httpServer.Serve())

}
func NewApplication(socketServer ports.SocketServer) app.Application {
	const FILE string = "/app/docker/scripts/websocket_database.sqlite"
	db, err := sql.Open("sqlite3", FILE)
	if err != nil {
		log.Fatalln(err)
	}
	cRepo := adapters.NewSqliteClientRepository(db)
	nRepo := adapters.NewSqliteNotificationRepostory(db, notification.NewDefaultMapper())

	return app.Application{
		Commands: app.Commands{
			SentBroadcast:               command.NewSendNotificationBroadcastHandler(socketServer),
			SentNotification:            command.NewSendNotificationHandler(socketServer, cRepo, nRepo),
			LogInAndGetOldNotifications: command.NewLogInAndGetOldNotificationsHandler(socketServer, cRepo, nRepo),
			UpdateNotificationStatus:    command.NewUpdateNotificationStatusHandler(socketServer, cRepo, nRepo),
		},
		Queries: app.Queries{
			ClientList:       query.NewClientListHandler(cRepo),
			NotificationList: query.NewNotificationListHandler(nRepo),
		},
	}

}
