package log

import (
	"testing"

	"github.com/qa-dev/universe/observer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
	log "github.com/Sirupsen/logrus"
	"bytes"
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
	m.t.Log("MockObserver.Notify called!")

	return args.Error(0)
}

func TestLog_Event(t *testing.T) {
	o := observer.NewObservable()
	l := Log{logger: log.New()}
	var b bytes.Buffer
	l.logger.Out = &b
	o.Register(l)

	o.NotifyEvent("Event")
	time.Sleep(200 * time.Millisecond)

	assert.Contains(t, b.String(), "Event")
}
