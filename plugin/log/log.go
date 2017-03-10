package log

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/observer"
)

type Log struct{}

func init() {
	l := Log{}
	observer.Obs.Register(l)
}

func (l Log) Event(v interface{}) error {
	log.Info(v)
	return nil
}

func (l Log) Subscribe(v interface{}) error {
	log.Info(v)
	return nil
}
