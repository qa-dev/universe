package log

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/subscribe"
)

type Log struct {
	logger *log.Logger
}

func init() {
	l := Log{logger: log.New()}
	plugins.Obs.Register(l)
}

func (l Log) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{"Log", "log"}
}

func (l Log) ProcessEvent(name string, eventData event.Event) {
	l.logger.Info(name, eventData)
}

func (l Log) Subscribe(name string, data subscribe.SubscribeData) {
	l.logger.Info(name, data)
}

func (l Log) Unsubscribe(name string, data subscribe.UnsubscribeData) {
	l.logger.Info(name, data)
}
