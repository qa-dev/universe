package log

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
)

type Log struct {
	logger *log.Logger
}

func NewLog() Log {
	return Log{logger: log.New()}
}

func (l Log) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		Name:    "Log",
		Tag:     "log",
		Version: 1,
	}
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

func (l Log) LoadSubscriber(data []byte) error {
	return nil
}

func (l Log) Loaded() {}
