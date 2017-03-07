package observer

import log "github.com/Sirupsen/logrus"

type Observer interface {
	Notify(v interface{}) error
}

type Observable struct {
	observers []Observer
}

func (o *Observable) Add(v Observer) {
	o.observers = append(o.observers, v)
}

func (o *Observable) Notify(v interface{}) {
	for _, ob := range o.observers {
		go notify(ob, v)
	}
}

func notify(ob Observer, v interface{}) {
	err := ob.Notify(v)
	if err != nil {
		log.Println(err)
	}
}
