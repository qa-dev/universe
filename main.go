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

	msession, err := mgo.Dial(cfg.Mongo.Host + ":" + cfg.Mongo.Port)
	if err != nil {
		panic(err)
	}
	msession.SetSyncTimeout(5 * time.Second)

	kpr := keeper.NewKeeper(msession)

	pluginStorage := plugins.NewPluginStorage()
	//TODO: keeper надо прокидывать не в плагин, а в subscribe и в dispatcher
	pluginStorage.Register(web.NewPluginWeb(kpr))
	pluginStorage.Register(logPlugin.NewLog())

	eventService := event.NewEventService(eventRmq)
	subscribeService := subscribe.NewSubscribeService(pluginStorage)
	dispatcherService := dispatcher.NewDispatcher(eventRmq, pluginStorage)
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
