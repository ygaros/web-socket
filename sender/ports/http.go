package ports

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sender/app"
	"sender/app/command"
	"sender/domain"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type HttpServer interface {
	Serve() error
	HandleNotification(w http.ResponseWriter, r *http.Request)
	ReadAllNotificationFromDb(w http.ResponseWriter, r *http.Request)
	ReadAllClientsFromDb(w http.ResponseWriter, r *http.Request)
}
type httpServer struct {
	app  app.Application
	url  string
	port int
}

func (s *httpServer) HandleNotification(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	request := domain.NotificationRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		log.Println("Failed to parse json data:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(request.Recipients) == 0 {
		log.Printf("Sending to all clients: %s\n", request.Message)
		err = s.app.Commands.SentBroadcast.Handle(r.Context(), command.SendNotificationBroadcast{
			Message: request.Message,
			//TODO IMPLEMENT GET USER FROM CONTEXT
			Sender: request.Sender,
		})
		if err != nil {
			log.Println("Failed to send data to clients:", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf("Sending to %d clients: %s\n", len(request.Recipients), request.Message)
		err = s.app.Commands.SentNotification.Handle(r.Context(), command.SendNotification{
			Message: request.Message,
			//TODO IMPLEMENT GET USER FROM CONTEXT
			Sender:     request.Sender,
			Recipients: request.Recipients,
		})
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
func (s *httpServer) ReadAllNotificationFromDb(w http.ResponseWriter, r *http.Request) {
	notifications, err := s.app.Queries.NotificationList.Handle(r.Context(), nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if marshalled, err := json.Marshal(notifications); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Write(marshalled)
	}
}
func (s *httpServer) ReadAllClientsFromDb(w http.ResponseWriter, r *http.Request) {
	clients, err := s.app.Queries.ClientList.Handle(r.Context(), nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if marshalled, err := json.Marshal(clients); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Write(marshalled)
	}
}
func (s *httpServer) Serve() error {
	log.Println("Starting rest-server...")
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/", s.HandleNotification)
	router.Get("/notifications", s.ReadAllNotificationFromDb)
	router.Get("/clients", s.ReadAllClientsFromDb)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.url, s.port), router)
}
func NewRestServer(url string, port int, app *app.Application) HttpServer {
	return &httpServer{
		url:  url,
		port: port,
		app:  *app,
	}
}
