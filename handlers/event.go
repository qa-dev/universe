package handlers

import (
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/qa-dev/universe/event"
)

type EventHandler struct {
	eventService EventPublisher
}

type EventPublisher interface {
	Publish(event.Event) error
}

func NewEventHandler(eventService EventPublisher) *EventHandler {
	return &EventHandler{eventService}
}

func (h *EventHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	eventName := req.RequestURI[utf8.RuneCountInString("/e"):]
	if len(eventName) == 0 {
		resp.Write([]byte("FAIL: BLANK EVENT NAME"))
		return
	}
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.Write([]byte("FAIL:" + err.Error()))
		return
	}

	e := event.Event{eventName, payload}
	err = h.eventService.Publish(e)
	if err == nil {
		resp.Write([]byte("OK!"))
	} else {
		resp.Write([]byte("FAIL:" + err.Error()))
	}

}
