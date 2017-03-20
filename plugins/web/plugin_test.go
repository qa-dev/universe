package web

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/keeper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
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

var kpr *keeper.Keeper
var database *mgo.Database

func init() {
	mongoUri := os.Getenv("MONGO_URI")
	session, err := mgo.Dial(mongoUri)
	if err != nil {
		panic("cant connect mongo")
	}
	database = session.DB("test_plugin_web")
	kpr = keeper.NewKeeper(session)
	kpr.SetCustomDatabaseName("test_plugin_web")
}

func TestNewPluginWeb(t *testing.T) {
	p := NewPluginWeb(kpr)
	assert.NotNil(t, p.keeper)
}

func TestPluginWeb_GetPluginInfo(t *testing.T) {
	p := NewPluginWeb(kpr)
	assert.Equal(t, "web", p.GetPluginInfo().Tag)
	assert.Equal(t, "Web", p.GetPluginInfo().Name)
}

func TestPluginWeb_Subscribe(t *testing.T) {
	database.DropDatabase()
	p := NewPluginWeb(kpr)
	inJson := []byte(`{"event_name": "test", "url": "hello"}`)
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	var result []SubscribeData
	err = kpr.GetSubscribers(PluginTag, &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
}

func TestPluginWeb_Subscribe_WrongInput(t *testing.T) {
	p := NewPluginWeb(kpr)
	inJson := []byte("{}")
	err := p.Subscribe(inJson)
	assert.Error(t, err)
}

func TestPluginWeb_Unsubscribe(t *testing.T) {
	database.DropDatabase()
	p := NewPluginWeb(kpr)
	inJson := []byte(`{"event_name": "test", "url": "hello"}`)
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	var result []SubscribeData
	err = kpr.GetSubscribers(PluginTag, &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	err = p.Unsubscribe(inJson)
	assert.NoError(t, err)
	result = nil
	err = kpr.GetSubscribers(PluginTag, &result)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))
}

func TestPluginWeb_Unsubscribe_WrongInput(t *testing.T) {
	p := NewPluginWeb(kpr)
	inJson := []byte("{}")
	err := p.Unsubscribe(inJson)
	assert.Error(t, err)
}

func TestPluginWeb_Unsubscribe_NonExistentSubscriber(t *testing.T) {
	database.DropDatabase()
	p := NewPluginWeb(kpr)
	subscribeJson := []byte(`{"event_name": "test", "url": "hello"}`)
	err := p.Subscribe(subscribeJson)
	assert.NoError(t, err)
	var result []SubscribeData
	err = kpr.GetSubscribers(PluginTag, &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	unsubscribeJson := []byte(`{"event_name": "test", "url": "bye"}`)
	err = p.Unsubscribe(unsubscribeJson)
	assert.Error(t, err)
	result = nil
	err = kpr.GetSubscribers(PluginTag, &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
}

func TestPluginWeb_ProcessEvent(t *testing.T) {
	database.DropDatabase()
	p := NewPluginWeb(kpr)
	expectedUrl := "test_url"
	expectedData := []byte(`{"hello": "world"}`)
	fakeClient := FakePostClient{t, expectedUrl, expectedData}
	p.client = fakeClient
	data := event.Event{Name: "test_event", Payload: expectedData}
	err := p.Subscribe([]byte(`{"event_name": "test_event", "url": "test_url"}`))
	assert.NoError(t, err)
	p.ProcessEvent(&data)
	time.Sleep(1 * time.Second)
}
