package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/qa-dev/universe/subscribe"
)

type SubscribeHandler struct {
	subscribeService *subscribe.SubscribeService
}

func NewSubscribeHandler(subscribeService *subscribe.SubscribeService) *SubscribeHandler {
	return &SubscribeHandler{subscribeService}
}

func (h *SubscribeHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var subscribeData subscribe.Subscribe
	err := decoder.Decode(&subscribeData)
	if err != nil {
		log.Println("Bad subscribe request")
		resp.Write([]byte("FAIL: BAD REQUEST"))
		return
	}
	defer req.Body.Close()

	err = h.subscribeService.ProcessSubscribe(subscribeData)
	if err == nil {
		resp.Write([]byte("OK"))
	} else {
		resp.Write([]byte("FAIL:" + err.Error()))
	}

}
