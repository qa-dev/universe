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
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

var database *mgo.Database

func init() {
	mongoUri := os.Getenv("MONGO_URI")
	session, err := mgo.Dial(mongoUri)
	if err != nil {
		panic("cant connect mongo")
	}
	database = session.DB("test_plugin_web")
	database.DropDatabase()
}

func TestNewPluginWeb(t *testing.T) {
	c := database.C("test_new_plugin_web")
	p := NewPluginWeb(c)
	assert.NotNil(t, p.collection)
}

func TestPluginWeb_GetPluginInfo(t *testing.T) {
	c := database.C("test_plugin_web_get_plugin_info")
	p := NewPluginWeb(c)
	assert.Equal(t, "web", p.GetPluginInfo().Tag)
	assert.Equal(t, "Web", p.GetPluginInfo().Name)
}

func TestPluginWeb_Subscribe(t *testing.T) {
	c := database.C("test_plugin_web_subscribe")
	p := NewPluginWeb(c)
	inJson := []byte(`{"event_name": "test", "url": "hello"}`)
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	resCount, err := c.Find(bson.M{"eventname": "test"}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, resCount)
}

func TestPluginWeb_Subscribe_WrongInput(t *testing.T) {
	c := database.C("test_plugin_web_wrong_input")
	p := NewPluginWeb(c)
	inJson := []byte("{}")
	err := p.Subscribe(inJson)
	assert.Error(t, err)
}

func TestPluginWeb_Unsubscribe(t *testing.T) {
	c := database.C("test_plugin_web_unsubscribe")
	p := NewPluginWeb(c)
	inJson := []byte(`{"event_name": "test", "url": "hello"}`)
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	resCount, err := c.Find(bson.M{"eventname": "test"}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, resCount)
	err = p.Unsubscribe(inJson)
	assert.NoError(t, err)
	resCount, err = c.Find(bson.M{"eventname": "test"}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, resCount)
}

func TestPluginWeb_Unsubscribe_WrongInput(t *testing.T) {
	c := database.C("test_plugin_web_unsubscribe_wrong_input")
	p := NewPluginWeb(c)
	inJson := []byte("{}")
	err := p.Unsubscribe(inJson)
	assert.Error(t, err)
}

func TestPluginWeb_Unsubscribe_NonExistentSubscriber(t *testing.T) {
	c := database.C("test_plugin_web_non_existent_subscriber")
	p := NewPluginWeb(c)
	subscribeJson := []byte(`{"event_name": "test", "url": "hello"}`)
	err := p.Subscribe(subscribeJson)
	assert.NoError(t, err)
	resCount, err := c.Find(bson.M{"eventname": "test"}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, resCount)
	unsubscribeJson := []byte(`{"event_name": "test", "url": "bye"}`)
	err = p.Unsubscribe(unsubscribeJson)
	assert.Error(t, err)
	resCount, err = c.Find(bson.M{"eventname": "test"}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, resCount)
}

func TestPluginWeb_ProcessEvent(t *testing.T) {
	c := database.C("test_plugin_web_process_event")
	p := NewPluginWeb(c)
	expectedUrl := "test_url"
	expectedData := []byte(`{"hello": "world"}`)
	fakeClient := FakePostClient{t, expectedUrl, expectedData}
	p.client = fakeClient
	data := event.Event{Name: "test_event", Payload: expectedData}
	err := p.Subscribe([]byte(`{"event_name": "test_event", "url": "test_url"}`))
	assert.NoError(t, err)
	p.ProcessEvent(data)
	time.Sleep(1 * time.Second)
}
