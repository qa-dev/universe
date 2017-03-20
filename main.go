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
	"github.com/qa-dev/universe/keeper"
	"github.com/qa-dev/universe/plugins"
	logPlugin "github.com/qa-dev/universe/plugins/log"
	"github.com/qa-dev/universe/plugins/web"
	"github.com/qa-dev/universe/queue"
	"github.com/qa-dev/universe/rabbitmq"
	"github.com/qa-dev/universe/subscribe"
	mgo "gopkg.in/mgo.v2"
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
	time.Sleep(2 * time.Second)
	defer eventRmq.Close()

	eventQueue := queue.NewQueue(eventRmq)

	msession, err := mgo.Dial(cfg.Mongo.Host + ":" + cfg.Mongo.Port)
	if err != nil {
		panic(err)
	}

	kpr := keeper.NewKeeper(msession)

	pluginStorage := plugins.NewPluginStorage()
	pluginStorage.Register(web.NewPluginWeb(kpr))
	pluginStorage.Register(logPlugin.NewLog())

	eventService := event.NewEventService(eventQueue)
	subscribeService := subscribe.NewSubscribeService(pluginStorage)
	dispatcherService := dispatcher.NewDispatcher(eventQueue, pluginStorage)
	dispatcherService.Run()

	for _, plg := range pluginStorage.GetPlugins() {
		plg.Loaded()
	}

	mux := http.NewServeMux()

	mux.Handle("/e/", handlers.NewEventHandler(eventService))
	mux.Handle("/subscribe/", handlers.NewSubscribeHandler(subscribeService))
	mux.Handle("/unsubscribe/", handlers.NewUnsubscribeHandler(subscribeService))

	listenData := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Info("Connected plugins:")
	for _, plg := range pluginStorage.GetPlugins() {
		log.Info(plg.GetPluginInfo().Name)
	}

	log.Info("App listen at ", listenData)
	log.Fatal(http.ListenAndServe(listenData, mux))
}
