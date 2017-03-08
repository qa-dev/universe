package main

import (
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/config"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/handlers"
	"github.com/qa-dev/universe/service"
	"github.com/qa-dev/universe/storage"
	"github.com/qa-dev/universe/subscribe"
)

func main() {
	cfg := config.LoadConfig()
	listenHost := cfg.GetString("app.host")
	listenPort := cfg.GetString("app.port")
	c := make(chan event.Event)
	storageUnit := storage.NewStorage()
	eventService := event.NewEventService(c)
	subscribeService := subscribe.NewSubscribeService(storageUnit)
	httpClient := &http.Client{}
	dispatcher := service.NewDispatcher(c, storageUnit, httpClient)
	dispatcher.Run()

	mux := http.NewServeMux()

	mux.Handle("/e/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe", handlers.NewSubscribeHandler(subscribeService))

	listenData := net.JoinHostPort(listenHost, listenPort)
	log.Info("App listen at ", listenData)
	log.Fatal(http.ListenAndServe(listenData, mux))
}
