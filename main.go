package main

import (
	"log"
	"net"
	"net/http"

	"github.com/qa-dev/universe/data"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/handlers"
	"github.com/qa-dev/universe/service"
	"github.com/qa-dev/universe/subscribe"
)

func main() {
	listenPort := "9713"
	c := make(chan event.Event)
	storage := data.NewStorage()
	eventService := event.NewEventService(c)
	subscribeService := subscribe.NewSubscribeService(storage)
	httpClient := &http.Client{}
	dispatcher := service.NewDispatcher(c, storage, httpClient)
	dispatcher.Run()

	mux := http.NewServeMux()

	mux.Handle("/event/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe", handlers.NewSubscribeHandler(subscribeService))

	log.Fatal(http.ListenAndServe(net.JoinHostPort("", listenPort), mux))
}
