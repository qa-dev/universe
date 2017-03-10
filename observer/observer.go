package observer

import log "github.com/Sirupsen/logrus"

var Obs *Observable

type Observer interface {
	Event(v interface{}) error
	Subscribe(v interface{}) error
}

type Observable struct {
	observers []Observer
}

func init() {
	Obs = NewObservable()
}

func NewObservable() *Observable {
	return &Observable{}
}

func (o *Observable) Register(v Observer) {
	o.observers = append(o.observers, v)
}

func (o *Observable) NotifyEvent(v interface{}) {
	for _, ob := range o.observers {
		go notify(ob, v)
	}
}

func notify(ob Observer, v interface{}) {
	err := ob.Event(v)
	if err != nil {
		log.Println(err)
	}
}
