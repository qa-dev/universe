package web

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/qa-dev/universe/event"
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

func (c FakePostClient) Do(r *http.Request) (*http.Response, error) {
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

func TestNewPluginWeb(t *testing.T) {
	p := NewPluginWeb()
	assert.NotNil(t, p.storage)
}

func TestPluginWeb_GetPluginInfo(t *testing.T) {
	p := NewPluginWeb()
	assert.Equal(t, "web", p.GetPluginInfo().Tag)
	assert.Equal(t, "Web", p.GetPluginInfo().Name)
}

func TestPluginWeb_Subscribe(t *testing.T) {
	p := NewPluginWeb()
	inJson := []byte("{\"event_name\": \"test\", \"url\": \"hello\"}")
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 1)
	assert.Len(t, p.storage.Data["test"], 1)
	assert.Equal(t, "hello", p.storage.Data["test"][0])
}

func TestPluginWeb_Subscribe_WrongInput(t *testing.T) {
	p := NewPluginWeb()
	inJson := []byte("{}")
	err := p.Subscribe(inJson)
	assert.Error(t, err)
}

func TestPluginWeb_Unsubscribe(t *testing.T) {
	p := NewPluginWeb()
	inJson := []byte("{\"event_name\": \"test\", \"url\": \"hello\"}")
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 1)
	err = p.Unsubscribe(inJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 0)
}

func TestPluginWeb_Unsubscribe_WrongInput(t *testing.T) {
	p := NewPluginWeb()
	inJson := []byte("{}")
	err := p.Unsubscribe(inJson)
	assert.Error(t, err)
}

func TestPluginWeb_Unsubscribe_NonExistentSubscriber(t *testing.T) {
	p := NewPluginWeb()
	subscribeJson := []byte("{\"event_name\": \"test\", \"url\": \"hello\"}")
	err := p.Subscribe(subscribeJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 1)
	unsubscribeJson := []byte("{\"event_name\": \"test\", \"url\": \"bye\"}")
	err = p.Unsubscribe(unsubscribeJson)
	assert.Error(t, err)
	assert.Len(t, p.storage.Data, 1)
}

func TestPluginWeb_ProcessEvent(t *testing.T) {
	p := NewPluginWeb()
	expectedUrl := "test_url"
	expectedData := []byte("{\"hello\": \"world\"}")
	fakeClient := FakePostClient{t, expectedUrl, expectedData}
	p.client = fakeClient
	data := event.Event{Name: "test_event", Payload: expectedData}
	err := p.Subscribe([]byte("{\"event_name\": \"test_event\", \"url\": \"test_url\"}"))
	assert.NoError(t, err)
	p.ProcessEvent(data)
}
