package rabbitmq

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

const QueueMaxPriority int16 = 255

type RabbitMQ struct {
	connection     *amqp.Connection
	channel        *amqp.Channel
	uri            string
	queueName      string
	closeError     chan *amqp.Error
	shouldBeClosed bool
}

func NewRabbitMQ(uri string, queueName string) *RabbitMQ {
	r := &RabbitMQ{uri: uri, queueName: queueName, shouldBeClosed: false}

	r.closeError = make(chan *amqp.Error)
	// Запускаем коннектор в горутине, который ждет ошибку коннекта и реконнектит
	go r.checkConnection()
	// Отправляем "ошибку" для первого запуска
	r.closeError <- amqp.ErrClosed

	return r
}

func (r *RabbitMQ) connect() {
	for {
		conn, err := amqp.Dial(r.uri)
		if err == nil {
			r.connection = conn
			return
		}
		log.Println(err)
		log.Printf("Trying to reconnect to RabbitMQ at %s\n", r.uri)
		time.Sleep(3000 * time.Millisecond)
	}
}

func (r *RabbitMQ) checkConnection() {
	var rabbitErr *amqp.Error
	for {
		rabbitErr = <-r.closeError
		if rabbitErr != nil {
			log.Printf("Connecting to %s\n", r.uri)

			r.connect()

			r.closeError = make(chan *amqp.Error)
			r.connection.NotifyClose(r.closeError)

			ch, err := r.connection.Channel()
			if err != nil {
				log.Println(err)
			}
			r.channel = ch
			r.DeclareQueue(r.queueName)
		}
	}
}

func (r *RabbitMQ) Close() {
	r.connection.Close()
	r.channel.Close()
}

// PublishWithPriority упаковывает msg в JSON и отправляет по ch каналу
// с ключем routingKey с приоритетом
func (r *RabbitMQ) PublishWithPriority(msg interface{}, priority uint8) error {
	return r.publishToEx("", r.queueName, msg, priority)
}

// PublishToEx упаковывает msg в JSON и отправляет по ch каналу
// с ключем routingKey в exchange
func (r *RabbitMQ) publishToEx(exchange, routingKey string, msg interface{}, priority uint8) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	//Publish message to RMQ
	err = r.channel.Publish(
		exchange,
		routingKey,
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
func (r *RabbitMQ) Consume(workerName string) (<-chan Delivery, error) {
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

func (r *RabbitMQ) DeclareQueue(queueName string) (*amqp.Queue, error) {
	args := make(amqp.Table)
	args["x-max-priority"] = QueueMaxPriority
	q, err := r.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
