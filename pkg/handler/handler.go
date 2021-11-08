package handler

import (
	"broker/pkg/broker"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func HandleConnections(responseWriter http.ResponseWriter, request *http.Request) {
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

	newClient := broker.MakeClient(rand.Intn(9999999), messageChan)

	broker.Broker().Clients = append(broker.Broker().Clients, newClient)

	defer func() {
		newClient.Disconnect()
	}()

	go func() {
		<-request.Context().Done()

		newClient.Disconnect()
	}()

	for {
		fmt.Fprintf(responseWriter, "data: %s\n\n", <-messageChan)

		flusher.Flush()
	}
}

func ReceiveEvent(response http.ResponseWriter, request *http.Request) {
	eventData := request.URL.Path
	broker.Broker().Notification <- eventData
	log.Print("Event received: " + eventData)
}
