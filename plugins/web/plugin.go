package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
)

type PluginWeb struct {
	storage *Storage
}

func NewPluginWeb() PluginWeb {
	storage := NewStorage()
	return PluginWeb{storage}
}

func init() {
	p := NewPluginWeb()
	plugins.Obs.Register(p)
}

func (p PluginWeb) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{Name: "Web", Tag: "web"}
}

func (p PluginWeb) Subscribe(input []byte) error {
	var subscribeData SubscribeData
	err := json.Unmarshal(input, &subscribeData)
	if err != nil {
		log.Errorf("%+v", input)
		errorText := fmt.Sprintf("Invalid input data: %s", string(input))
		return errors.New(errorText)
	}
	p.storage.AddSubscriber(subscribeData)
	return nil
}

func (p PluginWeb) Unsubscribe(input []byte) error {
	var unsubscribeData UnsubscribeData
	err := json.Unmarshal(input, &unsubscribeData)
	if err != nil {
		log.Errorf("%+v", input)
		errorText := fmt.Sprintf("Invalid input data: %s", string(input))
		return errors.New(errorText)
	}
	err = p.storage.RemoveSubscriber(unsubscribeData)

	return err
}

func (p PluginWeb) ProcessEvent(eventData event.Event) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(eventData.Payload))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for _, subscribeUrl := range p.storage.Data[eventData.Name] {
		req.URL, err = url.Parse(subscribeUrl)
		if err != nil {
			log.Error(err)
			continue
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Println("Status of sending event", eventData.Name, "is", resp.Status)
	}
}
