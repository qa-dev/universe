package rabbitmq

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

const QueueMaxPriority int16 = 255

type RabbitMQ struct {
	connection             *amqp.Connection
	channel                *amqp.Channel
	isOnline               bool
	errorCloseChan         chan *amqp.Error
	uri                    string
	queueName              string
	buffer                 chan interface{}
	reconnectWatchers      []chan *error
	reconnectWatchersMutex *sync.Mutex
}

func NewRabbitMQ(uri string, queueName string) *RabbitMQ {
	rabbit := RabbitMQ{
		isOnline:  false,
		uri:       uri,
		queueName: queueName,
	}
	rabbit.buffer = make(chan interface{})
	rabbit.errorCloseChan = make(chan *amqp.Error)
	rabbit.reconnectWatchers = make([]chan *error, 0)
	rabbit.reconnectWatchersMutex = &sync.Mutex{}
	go rabbit.monitorConnection()
	go rabbit.runBufferWorker()
	// first error for connect
	rabbit.errorCloseChan <- amqp.ErrClosed
	return &rabbit
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.connection.Close()
}

func (r *RabbitMQ) monitorConnection() {
	log.Info("Starting monitoring rabbitmq connection...")
	for {
		if err := <-r.errorCloseChan; err != nil {
			log.Info("Lost connection to rabbitmq.")
			log.Info(err)
			r.isOnline = false
			r.connect()
		}
	}
}

func (r *RabbitMQ) connect() {
	try := 0
	for {
		try += 1
		log.Infof("Reconnecting to rabbitmq... try %d...", try)
		conn, err := amqp.Dial(r.uri)
		if err != nil {
			log.Infof("Cant connect to amqp %s", r.uri)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		r.connection = conn
		r.errorCloseChan = make(chan *amqp.Error)
		r.connection.NotifyClose(r.errorCloseChan)
		ch, err := conn.Channel()
		if err != nil {
			log.Info(err)
			log.Info("Error getting channel, reconnecting.")
			r.connection.Close()
			continue
		}
		r.channel = ch
		_, err = r.DeclareQueue(r.queueName)
		if err != nil {
			log.Info(err)
			log.Info("Error declaring queue, reconnecting.")
			r.channel.Close()
			r.connection.Close()
			continue
		}

		r.sendReconnectNotifies(errors.New("reconnected"))
		r.isOnline = true
		log.Infof("Connection to rabbitmq established in %d tries.", try)
		return
	}
}

func (r *RabbitMQ) NotifyReconnect(receiver chan *error) {
	r.reconnectWatchersMutex.Lock()
	r.reconnectWatchers = append(r.reconnectWatchers, receiver)
	r.reconnectWatchersMutex.Unlock()
}

func (r *RabbitMQ) sendReconnectNotifies(err error) {
	r.reconnectWatchersMutex.Lock()
	observers := make([]chan *error, len(r.reconnectWatchers))
	copy(observers, r.reconnectWatchers)
	r.reconnectWatchersMutex.Unlock()
	for _, c := range observers {
		c <- &err
	}
}

func (r *RabbitMQ) DeclareQueue(name string) (*amqp.Queue, error) {
	args := make(amqp.Table)
	args["x-max-priority"] = QueueMaxPriority
	q, err := r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *RabbitMQ) IsOnline() bool {
	return r.isOnline
}

func (r *RabbitMQ) SendEvent(event interface{}) {
	r.buffer <- event
}

func (r *RabbitMQ) runBufferWorker() {
	for {
		ev := <-r.buffer
		// wait for rabbit online then immediately send event
		for r.IsOnline() == false {
			time.Sleep(500 * time.Millisecond)
		}
		r.sendEventToRabbit(ev, 1)
	}
}

func (r *RabbitMQ) sendEventToRabbit(msg interface{}, priority uint8) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	//Publish message to RMQ
	err = r.channel.Publish(
		"",
		r.queueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "encoding/json",
			Body:        body,
			Priority:    priority,
		})
	return err
}

// Вычитываем сообщения по одному из очереди.
func (r *RabbitMQ) GetConsumer(workerName string) (<-chan Delivery, error) {
	err := r.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	msgs, err := r.channel.Consume(
		r.queueName, // queue
		workerName,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arg s
	)
	if err != nil {
		return nil, err
	}
	return newDelivery(msgs), nil
}
