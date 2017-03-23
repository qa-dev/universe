package plugins

import (
	"testing"
	"time"

	"github.com/qa-dev/universe/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMsg = []byte("test")

type FakePlugin struct{}

func (FakePlugin) LoadSubscriber(data []byte) error {
	return nil
}

func (m FakePlugin) GetPluginInfo() *PluginInfo {
	return &PluginInfo{
		Name: "Name",
		Tag:  "fake",
	}
}

func (FakePlugin) Subscribe(input []byte) error {
	return nil
}

func (FakePlugin) Unsubscribe(input []byte) error {
	return nil
}

func (FakePlugin) Loaded()                             {}
func (FakePlugin) ProcessEvent(eventData *event.Event) {}

type MockObserver struct {
	a *assert.Assertions
	t *testing.T
	mock.Mock
}

func (MockObserver) LoadSubscriber(data []byte) error {
	return nil
}

func (m MockObserver) GetPluginInfo() *PluginInfo {
	return &PluginInfo{
		Name: "Name",
		Tag:  "name",
	}
}

func (m MockObserver) ProcessEvent(data *event.Event) {
	m.Called(data)
	m.a.Equal(string(testMsg), data.Name)
	m.t.Log("MockObserver.Event called!")
}

func (m MockObserver) Subscribe(input []byte) error {
	m.Called(input)
	m.a.Equal(testMsg, input)
	m.t.Log("MockObserver.Subscribe called!")
	return nil
}

func (m MockObserver) Unsubscribe(input []byte) error {
	m.Called(input)
	m.a.Equal(testMsg, input)
	m.t.Log("MockObserver.UnsubscribeData called!")
	return nil
}

func (m MockObserver) Loaded() {}

func TestObservable_Add(t *testing.T) {
	a := assert.New(t)

	o := PluginStorage{}
	ob1 := &MockObserver{a: a, t: t}

	ob1.On("ProcessEvent", &event.Event{string(testMsg), []byte(`{}`)}).Return(nil)
	o.Register(ob1)
	ev := event.Event{string(testMsg), []byte(`{}`)}
	o.ProcessEvent(&ev)
	time.Sleep(1 * time.Second)
}

func TestNewPluginStorage(t *testing.T) {
	storage := NewPluginStorage()
	assert.NotNil(t, storage)
}

func TestPluginStorage_ProcessSubscribe_WrongPluginName(t *testing.T) {
	storage := NewPluginStorage()
	err := storage.ProcessSubscribe("pew", []byte(""))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No plugin found")
}

func TestPluginStorage_ProcessSubscribe(t *testing.T) {
	storage := NewPluginStorage()
	storage.Register(FakePlugin{})
	err := storage.ProcessSubscribe("fake", []byte(""))
	assert.NoError(t, err)
}

func TestPluginStorage_GetPlugins(t *testing.T) {
	storage := NewPluginStorage()
	assert.Len(t, storage.GetPlugins(), 0)
	storage.Register(FakePlugin{})
	assert.Len(t, storage.GetPlugins(), 1)
}
