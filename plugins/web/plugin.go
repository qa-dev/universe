package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/keeper"
	"github.com/qa-dev/universe/plugins"
)

const PluginTag string = "web"

type PluginWeb struct {
	client HttpRequester
	keeper *keeper.Keeper
}

type HttpRequester interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewPluginWeb(kpr *keeper.Keeper) *PluginWeb {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return &PluginWeb{httpClient, kpr}
}

func (p *PluginWeb) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		Name:    "Web",
		Tag:     PluginTag,
		Version: 1,
	}
}

func (p *PluginWeb) Subscribe(input []byte) error {
	var subscribeData SubscribeData
	err := json.Unmarshal(input, &subscribeData)
	if err != nil || len(subscribeData.EventName) == 0 || len(subscribeData.Url) == 0 {
		log.Errorf("%+v", input)
		errorText := fmt.Sprintf("Invalid input data: %s", string(input))
		return errors.New(errorText)
	}
	return p.keeper.StoreSubscriber(PluginTag, subscribeData)
}

func (p *PluginWeb) Unsubscribe(input []byte) error {
	var unsubscribeData UnsubscribeData
	err := json.Unmarshal(input, &unsubscribeData)
	if err != nil || len(unsubscribeData.EventName) == 0 || len(unsubscribeData.Url) == 0 {
		log.Errorf("%+v", input)
		errorText := fmt.Sprintf("Invalid input data: %s", string(input))
		return errors.New(errorText)
	}
	return p.keeper.RemoveSubscriber(PluginTag, unsubscribeData)
}

func (p *PluginWeb) ProcessEvent(eventData event.Event) {
	var result []SubscribeData
	p.keeper.GetSubscribers(PluginTag, &result)
	for _, data := range result {
		go p.sendRequest(data.Url, eventData.Payload)
	}
}

func (p *PluginWeb) sendRequest(url string, payload []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: log statistics
	p.client.Do(req)
}

func (p *PluginWeb) Loaded() {}
