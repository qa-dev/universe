package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMsg = "test"

type MockObserver struct {
	a *assert.Assertions
	t *testing.T
	mock.Mock
}

func (m MockObserver) Event(v interface{}) error {
	args := m.Called(v)
	m.a.Equal(testMsg, v.(string))
	m.t.Log("MockObserver.Event called!")

	return args.Error(0)
}

func (m MockObserver) Subscribe(v interface{}) error {
	args := m.Called(v)
	m.a.Equal(testMsg, v.(string))
	m.t.Log("MockObserver.Subscribe called!")

	return args.Error(0)
}

func TestObservable_Add(t *testing.T) {
	a := assert.New(t)

	o := Observable{}
	ob1 := &MockObserver{a: a, t: t}

	ob1.On("Notify", testMsg).Return(nil)
	o.Register(ob1)
	o.NotifyEvent(testMsg)
}
