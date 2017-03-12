package log

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
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

func (l Log) ProcessEvent(eventData event.Event) {
	l.logger.Info(eventData)
}

func (l Log) Subscribe(input []byte) error {
	l.logger.Info(string(input))
	return nil
}

func (l Log) Unsubscribe(input []byte) error {
	l.logger.Info(string(input))
	return nil
}
