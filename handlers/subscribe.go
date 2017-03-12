package handlers

import (
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/qa-dev/universe/subscribe"
)

type SubscribeHandler struct {
	subscribeService *subscribe.SubscribeService
}

func NewSubscribeHandler(subscribeService *subscribe.SubscribeService) *SubscribeHandler {
	return &SubscribeHandler{subscribeService}
}

func (h *SubscribeHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	pluginName := req.RequestURI[utf8.RuneCountInString("/subscribe/"):]
	if len(pluginName) == 0 {
		resp.Write([]byte("FAIL: BLANK PLUGIN NAME"))
		return
	}
	input, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.Write([]byte("FAIL: BAD REQUEST"))
		return
	}
	defer req.Body.Close()

	err = h.subscribeService.ProcessSubscribe(pluginName, input)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "` + err.Error() + `"}"`))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"status": "ok"}`))

}
