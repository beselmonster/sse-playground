package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

var messagesBroker *MessagesBroker = &MessagesBroker{
	clients:      make([]Client, 0),
	notification: make(chan string, 1),
}

type MessagesBroker struct {
	clients []Client

	notification chan string
}

//Listen for events
func (messagesBroker *MessagesBroker) listenForNewEvents() {
	for {
		event := <-messagesBroker.notification

		for _, client := range messagesBroker.clients {
			client.connection <- event
		}
	}
}

func main() {
	go messagesBroker.listenForNewEvents()

	http.HandleFunc("/events/listen/", eventHandler)

	http.HandleFunc("/events/push/", receiveEvent)

	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}

func eventHandler(responseWriter http.ResponseWriter, request *http.Request) {
	flusher, ok := responseWriter.(http.Flusher)

	if !ok {
		http.Error(responseWriter, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	log.Print("Client connected")

	messageChan := make(chan string)

	newClient := Client{id: rand.Int(), connection: messageChan}

	messagesBroker.clients = append(messagesBroker.clients, newClient)

	defer func() {
		newClient.disconnect()
	}()

	go func() {
		<-request.Context().Done()

		newClient.disconnect()
	}()

	for {
		fmt.Fprintf(responseWriter, "data: %s\n\n", <-messageChan)

		flusher.Flush()
	}
}

func receiveEvent(response http.ResponseWriter, request *http.Request) {
	eventData := request.URL.Path
	messagesBroker.notification <- eventData
	log.Print("Event received: " + eventData)
}

type Client struct {
	id int

	connection chan string
}

func (clientToDisconnect *Client) disconnect() {
	clients := make([]Client, 0)
	for _, client := range messagesBroker.clients {
		if client.id != clientToDisconnect.id {
			clients = append(clients, client)
		}
	}
	messagesBroker.clients = clients

	log.Print("Connection closed with id ", clientToDisconnect.id)
}
