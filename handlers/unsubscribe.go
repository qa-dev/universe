package handlers

import (
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/qa-dev/universe/subscribe"
)

type UnsubscribeHandler struct {
	SubscribeService *subscribe.SubscribeService
}

func NewUnsubscribeHandler(subscribeService *subscribe.SubscribeService) *UnsubscribeHandler {
	return &UnsubscribeHandler{subscribeService}
}

func (h *UnsubscribeHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	pluginName := req.RequestURI[utf8.RuneCountInString("/unsubscribe/"):]
	if len(pluginName) == 0 {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(`{"error": "Blank plugin name"}`))
		return
	}
	input, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(`{"error": "Bad request"}`))
		return
	}
	defer req.Body.Close()

	err = h.SubscribeService.ProcessUnsubscribe(pluginName, input)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"status": "ok"}`))

}
