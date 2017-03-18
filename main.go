package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/config"
	"github.com/qa-dev/universe/dispatcher"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/handlers"
	"github.com/qa-dev/universe/plugins"
	logPlugin "github.com/qa-dev/universe/plugins/log"
	"github.com/qa-dev/universe/plugins/web"
	"github.com/qa-dev/universe/rabbitmq"
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

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	eventRmq := rabbitmq.NewRabbitMQ(cfg.Rmq.Uri, cfg.Rmq.EventQueue)
	time.Sleep(2 * time.Second)
	defer eventRmq.Close()

	pluginStorage := plugins.NewPluginStorage()
	pluginStorage.Register(web.NewPluginWeb())
	pluginStorage.Register(logPlugin.NewLog())

	eventService := event.NewEventService(eventRmq)
	subscribeService := subscribe.NewSubscribeService(pluginStorage)
	dispatcherService := dispatcher.NewDispatcher(eventRmq, pluginStorage)
	dispatcherService.Run()

	mux := http.NewServeMux()

	mux.Handle("/e/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe/", handlers.NewSubscribeHandler(subscribeService))

	listenData := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Info("Connected plugins:")
	for _, plg := range pluginStorage.GetPlugins() {
		log.Info(plg.GetPluginInfo().Name)
	}

	srv := &http.Server{Addr: listenData, Handler: mux}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
		log.Info("App listen at ", listenData)
	}()

	<-stopChan
	log.Info("Shutting down server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Info("Server gracefully stopped")
}
