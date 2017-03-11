package log

import (
	"bytes"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/stretchr/testify/assert"
)

func TestLog_ProcessEvent(t *testing.T) {
	o := plugins.NewObservable()
	l := Log{logger: log.New()}
	var b bytes.Buffer
	l.logger.Out = &b
	o.Register(l)

	o.ProcessEvent(event.Event{"Event", []byte(`{\"hello\": \"test\"}`)})
	time.Sleep(200 * time.Millisecond)

	t.Log(b.String())
	assert.Contains(t, b.String(), "Event")
}
