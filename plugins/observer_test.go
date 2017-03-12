package plugins

import (
	"testing"

	"time"

	"github.com/qa-dev/universe/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMsg = []byte("test")

type MockObserver struct {
	a *assert.Assertions
	t *testing.T
	mock.Mock
}

func (m MockObserver) GetPluginInfo() *PluginInfo {
	return &PluginInfo{"Name", "name"}
}

func (m MockObserver) ProcessEvent(data event.Event) {
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

func TestObservable_Add(t *testing.T) {
	a := assert.New(t)

	o := Observable{}
	ob1 := &MockObserver{a: a, t: t}

	ob1.On("ProcessEvent", event.Event{string(testMsg), []byte(`{}`)}).Return(nil)
	o.Register(ob1)
	o.ProcessEvent(event.Event{string(testMsg), []byte(`{}`)})
	time.Sleep(1 * time.Second)
}
