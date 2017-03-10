package log

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/observer"
)

type Log struct{
	logger *log.Logger
}

func init() {
	l := Log{logger: log.New()}
	observer.Obs.Register(l)
}

func (l Log) Event(v interface{}) error {
	l.logger.Info(v)
	return nil
}

func (l Log) Subscribe(v interface{}) error {
	l.logger.Info(v)
	return nil
}
