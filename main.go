package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/config"
	"github.com/qa-dev/universe/dispatcher"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/handlers"
	"github.com/qa-dev/universe/plugins"
	logPlugin "github.com/qa-dev/universe/plugins/log"
	"github.com/qa-dev/universe/plugins/web"
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

	//eventRmq := rabbitmq.NewRabbitMQ(cfg.Rmq.Uri, cfg.Rmq.EventQueue)
	//time.Sleep(2 * time.Second)
	//defer eventRmq.Close()

	queue := event.NewEventQueue()

	pluginStorage := plugins.NewPluginStorage()
	pluginStorage.Register(web.NewPluginWeb())
	pluginStorage.Register(logPlugin.NewLog())

	eventService := event.NewEventService(queue)
	subscribeService := subscribe.NewSubscribeService(pluginStorage)
	dispatcherService := dispatcher.NewDispatcher(queue, pluginStorage)
	dispatcherService.Run()

	mux := http.NewServeMux()

	mux.Handle("/e/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe/", handlers.NewSubscribeHandler(subscribeService))

	listenData := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Info("Connected plugins:")
	for _, plg := range pluginStorage.GetPlugins() {
		log.Info(plg.GetPluginInfo().Name)
	}
	log.Info("App listen at ", listenData)
	log.Fatal(http.ListenAndServe(listenData, mux))
}
