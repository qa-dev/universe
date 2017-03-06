package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/qa-dev/Universe/data"
	"github.com/qa-dev/Universe/service"
)

type SubscribeHandler struct {
	subscribeService *service.SubscribeService
}

func NewSubscribeHandler(subscribeService *service.SubscribeService) *SubscribeHandler {
	return &SubscribeHandler{subscribeService}
}

func (h *SubscribeHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var subscribe data.Subscribe
	err := decoder.Decode(&subscribe)
	if err != nil {
		log.Println("Bad subscribe request")
		resp.Write([]byte("FAIL: BAD REQUEST"))
		return
	}
	defer req.Body.Close()

	err = h.subscribeService.ProcessSubscribe(subscribe)
	if err == nil {
		resp.Write([]byte("OK"))
	} else {
		resp.Write([]byte("FAIL:" + err.Error()))
	}

}
