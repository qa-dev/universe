package rabbitmq

import (
	"github.com/streadway/amqp"
)

type DeliveryManager struct {
	inChan  <-chan amqp.Delivery
	outChan chan Delivery
}

type Delivery struct {
	msg amqp.Delivery
}

func newDelivery(ch <-chan amqp.Delivery) <-chan Delivery {
	outChan := make(chan Delivery)
	d := &DeliveryManager{inChan: ch, outChan: outChan}
	go d.retranslate()
	return (<-chan Delivery)(d.outChan)
}

func (dm *DeliveryManager) retranslate() {
	for {
		select {
		case msg := <-dm.inChan:
			out := Delivery{msg: msg}
			dm.outChan <- out
		}
	}
}

func (d Delivery) Ack() error {
	return d.msg.Ack(false)
}

func (d Delivery) Reject() error {
	return d.msg.Reject(true)
}

func (d Delivery) Body() []byte {
	return d.msg.Body
}

func (d *Delivery) Priority() uint8 {
	return d.msg.Priority
}
