package dispatcher

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/subscribe"
	"github.com/stretchr/testify/assert"
)

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
	queue := event.NewEventQueue()
	storage := plugins.NewPluginStorage()
	dsp := NewDispatcher(queue, storage)
	assert.Equal(t, fmt.Sprintf("%p", queue), fmt.Sprintf("%p", dsp.queue))
}

func TestDispatcher_Run(t *testing.T) {
	queue := event.NewEventQueue()
	storage := plugins.NewPluginStorage()
	requestData := []byte(`{"test": "test"}`)
	subscrService := subscribe.NewSubscribeService(storage)
	eventService := event.NewEventService(queue)
	subscribeData := []byte(`{"test": "hello"}`)
	subscrService.ProcessSubscribe("log", subscribeData)
	dsp := NewDispatcher(queue, storage)
	assert.NotNil(t, dsp)
	dsp.Run()
	err := eventService.Publish(event.Event{"test.event", requestData})
	assert.NoError(t, err)
	// TODO: assert log
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
