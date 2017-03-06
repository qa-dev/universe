package service

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/qa-dev/universe/data"
	"github.com/stretchr/testify/assert"
	"time"
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
	ch := make(chan event.Event)
	storage := data.NewStorage()
	client := &http.Client{}
	dsp := NewDispatcher(ch, storage, client)
	assert.Equal(t, fmt.Sprintf("%p", ch), fmt.Sprintf("%p", dsp.ch))
	assert.Equal(t, fmt.Sprintf("%p", storage), fmt.Sprintf("%p", dsp.storage))
	assert.Equal(t, fmt.Sprintf("%p", client), fmt.Sprintf("%p", dsp.httpClient))
}

func TestDispatcher_Run(t *testing.T) {
	requestUrl := "test_url"
	requestData := []byte("{\"test\": \"test\"}")
	ch := make(chan event.Event)
	storage := data.NewStorage()
	subscrService := NewSubscribeService(storage)
	eventService := NewEventService(ch)
	subscribe := data.Subscribe{EventName: "test.event", WebHookPath: requestUrl}
	subscrService.ProcessSubscribe(subscribe)
	client := &FakePostClient{t, requestUrl, requestData}
	dsp := NewDispatcher(ch, storage, client)
	assert.NotNil(t, dsp)
	dsp.Run()
	err := eventService.Publish("test.event", requestData)
	assert.NoError(t, err)
}

func TestDispatcher_Run_Negative(t *testing.T) {
	requestUrl := "test_url"
	requestData := []byte("{\"test\": \"test\"}")
	ch := make(chan event.Event)
	storage := data.NewStorage()
	subscrService := NewSubscribeService(storage)
	eventService := NewEventService(ch)
	subscribe := data.Subscribe{EventName: "test.event", WebHookPath: requestUrl}
	subscrService.ProcessSubscribe(subscribe)
	client := &FakePostClient{t, requestUrl, requestData}
	dsp := NewDispatcher(ch, storage, client)
	assert.NotNil(t, dsp)
	dsp.Run()
	err := eventService.Publish("test.event", requestData)
	assert.NoError(t, err)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	event := event.Event{Name: "", Payload: []byte("")}
	ch <- event

	log.SetOutput(os.Stderr)
	logVal := buf.String()

	time.Sleep(1000)

	assert.Contains(t, logVal, "No subscribers for event")
}
