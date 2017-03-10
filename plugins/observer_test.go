package plugins

import (
	"testing"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/subscribe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMsg = "test"

type MockObserver struct {
	a *assert.Assertions
	t *testing.T
	mock.Mock
}

func (m MockObserver) GetPluginInfo() *PluginInfo {
	return &PluginInfo{"Name", "name"}
}

func (m MockObserver) ProcessEvent(name string, data event.Event) {
	m.Called(name, data)
	m.a.Equal(testMsg, name)
	m.t.Log("MockObserver.Event called!")
}

func (m MockObserver) Subscribe(name string, data subscribe.SubscribeData) {
	m.Called(name, data)
	m.a.Equal(testMsg, name)
	m.t.Log("MockObserver.Subscribe called!")
}

func (m MockObserver) Unsubscribe(name string, data subscribe.UnsubscribeData) {
	m.Called(name, data)
	m.a.Equal(testMsg, name)
	m.t.Log("MockObserver.UnsubscribeData called!")
}

func TestObservable_Add(t *testing.T) {
	a := assert.New(t)

	o := Observable{}
	ob1 := &MockObserver{a: a, t: t}

	ob1.On("Notify", testMsg).Return(nil)
	o.Register(ob1)
	o.ProcessEvent(testMsg, event.Event{testMsg, []byte(`{}`)})
}
