package broker

import (
	"log"
)

type Client struct {
	id         int
	connection chan string
}

type MessagesBroker struct {
	Clients      []Client
	Notification chan string
}

var msgsBroker *MessagesBroker = &MessagesBroker{
	Clients:      make([]Client, 0),
	Notification: make(chan string, 1),
}

func Broker() *MessagesBroker {
	return msgsBroker
}

func MakeClient(id int, connection chan string) Client {
	return Client{id: id, connection: connection}
}

func (msgsBroker *MessagesBroker) ListenForNewEvents() {
	for {
		event := <-msgsBroker.Notification

		for _, client := range msgsBroker.Clients {
			client.connection <- event
		}
	}
}

func (clientToDisconnect *Client) Disconnect() {
	clients := make([]Client, 0)
	for _, client := range msgsBroker.Clients {
		if client.id != clientToDisconnect.id {
			clients = append(clients, client)
			break
		}
	}
	msgsBroker.Clients = clients

	log.Print("Connection closed with id ", clientToDisconnect.id)
}
