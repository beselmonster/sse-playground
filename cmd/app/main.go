package main

import (
	"broker/pkg/broker"
	"broker/pkg/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	port := "8081"

	log.Printf("SSE SERVER IS STARTING ON PORT %s ....", port)

	go broker.Broker().ListenForNewEvents()

	http.HandleFunc("/events/listen/", handler.HandleConnections)

	http.HandleFunc("/events/push/", handler.ReceiveEvent)

	log.Fatal("HTTP server error: ", http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}
