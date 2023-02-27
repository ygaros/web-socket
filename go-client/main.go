package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Notification struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
type NotificationStatusUpdate struct {
	Id         string `json:"id"`
	ClientName string `json:"clientName"`
	Status     string `json:"status"`
}

func main() {
	clientName := "goclient"
	connection, err := net.Dial("tcp", fmt.Sprintf(":%d", 9876))
	if err != nil {
		log.Fatalln(err)
	}
	defer connection.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for {
		login := "login:!" + clientName + "!"
		_, err := connection.Write([]byte(login))
		if err != nil {
			log.Fatalln(err)
		} else {
			break
		}
	}
	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		bufferToUnMarshal := buffer[:mLen]
		log.Printf("Received: %s\n", string(bufferToUnMarshal))
		notification := Notification{}
		notifications := make([]Notification, 0)
		if err := json.Unmarshal(bufferToUnMarshal, &notification); err == nil {
			nsu := NotificationStatusUpdate{
				Id:         notification.Id,
				ClientName: clientName,
				Status:     "received",
			}
			if marshaled, err := json.Marshal(nsu); err == nil {
				log.Printf("Sending update for notification %s to status %s\n", nsu.Id, nsu.Status)
				connection.Write(marshaled)
			}
		} else if err := json.Unmarshal(bufferToUnMarshal, &notifications); err == nil {
			parsed := make([]NotificationStatusUpdate, 0)
			for _, n := range notifications {
				nsu := NotificationStatusUpdate{
					Id:         n.Id,
					ClientName: clientName,
					Status:     "seen",
				}
				parsed = append(parsed, nsu)
			}
			if marshaled, err := json.Marshal(parsed); err == nil {
				log.Printf("Sending update for notifications len(%d) to status %s\n", len(parsed), "seen")
				connection.Write(marshaled)
			} else {
				log.Println(err)
			}
		}
	}
}
