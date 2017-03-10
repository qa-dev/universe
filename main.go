package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/config"
	"github.com/qa-dev/universe/dispatcher"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/handlers"
	_ "github.com/qa-dev/universe/plugin/log"
	"github.com/qa-dev/universe/rabbitmq"
	"github.com/qa-dev/universe/storage"
	"github.com/qa-dev/universe/subscribe"
)

func main() {
	cfgFile := flag.String("config", "./config.json", "Config file path")
	flag.Parse()

	cfg := &config.Config{}
	err := config.LoadFromFile(*cfgFile, cfg)
	if err != nil {
		log.Fatal(err)
	}

	eventRmq := rabbitmq.NewRabbitMQ(cfg.Rmq.Uri, cfg.Rmq.EventQueue)
	time.Sleep(5 * time.Second)
	defer eventRmq.Close()

	storageUnit := storage.NewStorage()
	eventService := event.NewEventService(eventRmq)
	subscribeService := subscribe.NewSubscribeService(storageUnit)
	httpClient := &http.Client{}
	dispatcherService := dispatcher.NewDispatcher(eventRmq, storageUnit, httpClient)
	dispatcherService.Run()

	mux := http.NewServeMux()

	mux.Handle("/e/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe/", handlers.NewSubscribeHandler(subscribeService))

	listenData := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Info("App listen at ", listenData)

	log.Fatal(http.ListenAndServe(listenData, mux))
}
