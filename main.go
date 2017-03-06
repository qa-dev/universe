package main

import (
	"log"
	"net"
	"net/http"

	"github.com/qa-dev/universe/data"
	"github.com/qa-dev/universe/handlers"
	"github.com/qa-dev/universe/service"
)

func main() {
	listenPort := "9713"
	c := make(chan data.Event)
	storage := data.NewStorage()
	subscribeService := service.NewSubscribeService(storage)
	eventService := service.NewEventService(c)
	httpClient := &http.Client{}
	dispatcher := service.NewDispatcher(c, storage, httpClient)
	dispatcher.Run()

	mux := http.NewServeMux()

	mux.Handle("/event/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe", handlers.NewSubscribeHandler(subscribeService))

	log.Fatal(http.ListenAndServe(net.JoinHostPort("", listenPort), mux))
}
