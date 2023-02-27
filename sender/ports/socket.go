package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"regexp"
	"sender/app"
	"sender/app/command"
	"sender/domain"
	"sender/domain/notification"
	"sync"

	"github.com/google/uuid"
)

const LOGIN_PATTERN = "login:![a-zA-Z]{3,}!"

type SocketServer interface {
	SendToAllClients(message domain.NotificationToSendBroadcast) error
	SendToRecipient(message domain.NotificationToSend, recipient string) error
	SendMultipleToRecipient(messages []domain.NotificationToSend, recipient string) error
	Serve() error
	SetApp(app *app.Application)
}
type socketServer struct {
	app app.Application

	lock                sync.RWMutex
	activeConnections   map[string]net.Conn
	url                 string
	port                int
	loginChannel        chan string
	updateStatusChannel chan domain.NotificationToChangeStatus
}

func (s *socketServer) SendToAllClients(message domain.NotificationToSendBroadcast) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, conn := range s.activeConnections {
		s.send(conn, domain.NotificationToSend{
			Id:      uuid.Nil,
			Message: message.Message,
			//TODO GET USER FROM CONTEXT
			Status: notification.BROADCAST.String(),
		})
	}
	return nil
}
func (s *socketServer) SendToRecipient(message domain.NotificationToSend, recipient string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	conn, ok := s.activeConnections[recipient]
	if ok {
		_, err := s.send(conn, message)
		return err
	}
	return fmt.Errorf("recipient %s isnt logged in", recipient)
}
func (s *socketServer) SendMultipleToRecipient(messages []domain.NotificationToSend, recipient string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	conn, ok := s.activeConnections[recipient]
	if ok {
		_, err := s.sendAll(conn, messages)
		return err
	}
	return fmt.Errorf("recipient %s isnt logged in", recipient)
}
func (s *socketServer) sendAll(conn net.Conn, messages []domain.NotificationToSend) (int, error) {
	if marshaled, err := json.Marshal(messages); err == nil {
		return conn.Write(marshaled)
	} else {
		return 0, err
	}
}
func (s *socketServer) send(conn net.Conn, message domain.NotificationToSend) (int, error) {
	if marshaled, err := json.Marshal(message); err == nil {
		return conn.Write(marshaled)
	} else {
		return 0, err
	}
}
func (s *socketServer) Serve() error {
	log.Println("Starting socket-server...")
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.url, s.port))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer server.Close()
	for {
		connection, err := server.Accept()
		if err != nil {
			log.Println(err)
			return nil
		}
		go s.handleClient(connection)
	}
}

func (s *socketServer) remove(recipient string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.activeConnections, recipient)
}

func (s *socketServer) handleClient(connection net.Conn) {
	reg, err := regexp.Compile(LOGIN_PATTERN)
	if err != nil {
		log.Println(err)
	}
	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		bufferToConvert := buffer[:mLen]
		message := string(bufferToConvert)
		if reg.Match([]byte(message)) {
			userName := retrieveUserName(message)
			log.Println(fmt.Sprintf("User %s logged in", userName))

			s.lock.Lock()
			s.activeConnections[userName] = connection
			s.lock.Unlock()
			// err := s.app.Commands.LogInAndGetOldNotifications.Handle(context.Background(), userName)
			// if err != nil {
			// 	continue
			// }
			go s.app.Commands.LogInAndGetOldNotifications.Handle(context.Background(), userName)
		} else {
			ntcs := domain.NotificationToChangeStatus{}
			ntcsMultiple := make([]domain.NotificationToChangeStatus, 0)
			log.Println(string(bufferToConvert))
			if err := json.Unmarshal(bufferToConvert, &ntcs); err == nil {
				status, err := notification.NewNotificationStatusFromString(ntcs.Status)
				if err != nil {
					log.Println(err)
					continue
				}
				// err = s.app.Commands.UpdateNotificationStatus.Handle(context.Background(), command.UpdateNotification{
				// 	Id:         ntcs.Id,
				// 	ClientName: ntcs.ClientName,
				// 	Status:     status,
				// })
				// if err != nil {
				// 	continue
				// }
				go s.app.Commands.UpdateNotificationStatus.Handle(context.Background(), command.UpdateNotification{
					Id:         ntcs.Id,
					ClientName: ntcs.ClientName,
					Status:     status,
				})
			} else if err := json.Unmarshal(bufferToConvert, &ntcsMultiple); err == nil {
				for _, ntcsFromMultiple := range ntcsMultiple {
					status, err := notification.NewNotificationStatusFromString(ntcsFromMultiple.Status)
					if err != nil {
						log.Println(err)
						continue
					}
					go s.app.Commands.UpdateNotificationStatus.Handle(context.Background(), command.UpdateNotification{
						Id:         ntcsFromMultiple.Id,
						ClientName: ntcsFromMultiple.ClientName,
						Status:     status,
					})
				}
			}
		}
	}
}
func (s *socketServer) SetApp(app *app.Application) {
	s.app = *app
}

func retrieveUserName(loginMessage string) string {
	return loginMessage[7 : len(loginMessage)-1]
}
func NewSocketServer(url string, port int) SocketServer {
	return &socketServer{
		activeConnections: make(map[string]net.Conn),
		url:               url,
		port:              port,
	}
}
