package service

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"time"

	"github.com/qa-dev/universe/config"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/rabbitmq"
	"github.com/qa-dev/universe/storage"
	"github.com/qa-dev/universe/subscribe"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.SetTestDitectory()
}

type FakeClosingBuffer struct {
	*bytes.Buffer
}

func (cb *FakeClosingBuffer) Close() error {
	return nil
}

type FakePostClient struct {
	t                   *testing.T
	ExpectedRequestUrl  string
	ExpectedRequestData []byte
}

func (c *FakePostClient) Do(r *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(c.t, err)
	assert.Equal(c.t, c.ExpectedRequestUrl, r.URL.String())
	assert.Equal(c.t, c.ExpectedRequestData, body)
	closingBuffer := &FakeClosingBuffer{bytes.NewBufferString("Hi!")}
	var readCloser io.ReadCloser
	readCloser = closingBuffer
	response := &http.Response{}
	response.Status = "200 OK"
	response.StatusCode = 200
	response.Body = readCloser
	return response, nil
}

func TestNewDispatcher(t *testing.T) {
	rmq := rabbitmq.NewRabbitMQ(config.LoadConfig().GetString("rmq.uri"), "test_event_service_push_event_queue")
	defer rmq.Close()
	// Даем время на подключение
	time.Sleep(5 * time.Second)
	storageUnit := storage.NewStorage()
	client := &http.Client{}
	dsp := NewDispatcher(rmq, storageUnit, client)
	assert.Equal(t, fmt.Sprintf("%p", rmq), fmt.Sprintf("%p", dsp.rmq))
	assert.Equal(t, fmt.Sprintf("%p", storageUnit), fmt.Sprintf("%p", dsp.storage))
	assert.Equal(t, fmt.Sprintf("%p", client), fmt.Sprintf("%p", dsp.httpClient))
}

func TestDispatcher_Run(t *testing.T) {
	rmq := rabbitmq.NewRabbitMQ(config.LoadConfig().GetString("rmq.uri"), "test_event_service_push_event_queue")
	defer rmq.Close()
	// Даем время на подключение
	time.Sleep(5 * time.Second)
	requestUrl := "test_url"
	requestData := []byte("{\"test\": \"test\"}")
	storageUnit := storage.NewStorage()
	subscrService := subscribe.NewSubscribeService(storageUnit)
	eventService := event.NewEventService(rmq)
	subscribeData := subscribe.Subscribe{EventName: "test.event", WebHookPath: requestUrl}
	subscrService.ProcessSubscribe(subscribeData)
	client := &FakePostClient{t, requestUrl, requestData}
	dsp := NewDispatcher(rmq, storageUnit, client)
	assert.NotNil(t, dsp)
	dsp.Run()
	err := eventService.Publish(event.Event{"test.event", requestData})
	assert.NoError(t, err)
}

// Disabled
//
//func TestDispatcher_Run_Negative(t *testing.T) {
//	requestUrl := "test_url"
//	requestData := []byte("{\"test\": \"test\"}")
//	ch := make(chan data.Event)
//	storage := data.NewStorage()
//	subscrService := subscribe.NewSubscribeService(storage)
//	eventService := NewEventService(ch)
//	subscribeData := subscribe.Subscribe{EventName: "test.event", WebHookPath: requestUrl}
//	subscrService.ProcessSubscribe(subscribeData)
//	client := &FakePostClient{t, requestUrl, requestData}
//	dsp := NewDispatcher(ch, storage, client)
//	assert.NotNil(t, dsp)
//	dsp.Run()
//	err := eventService.PushEvent("test.event", requestData)
//	assert.NoError(t, err)
//
//	var buf bytes.Buffer
//	log.SetOutput(&buf)
//
//	event := data.Event{Name: "", Payload: []byte("")}
//	ch <- event
//
//	log.SetOutput(os.Stderr)
//	logVal := buf.String()
//
//	time.Sleep(1000)
//
//	assert.Contains(t, logVal, "No subscribers for event")
//}
