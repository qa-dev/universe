package handlers

import (
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/qa-dev/Universe/service"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService}
}

func (h *EventHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	eventName := req.RequestURI[utf8.RuneCountInString("/event/"):]
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	err = h.eventService.PushEvent(eventName, payload)
	if err == nil {
		resp.Write([]byte("OK"))
	} else {
		resp.Write([]byte("FAIL:" + err.Error()))
	}

}
