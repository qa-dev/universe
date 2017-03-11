package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	p.storage.Mutex.Lock()
	if _, ok := p.storage.Data[subscribeData.EventName]; !ok {
		p.storage.Data[subscribeData.EventName] = make([]string, 0)
	}
	p.storage.Data[subscribeData.EventName] = append(p.storage.Data[subscribeData.EventName], subscribeData.Url)
	p.storage.Mutex.Unlock()

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
	p.storage.Mutex.Lock()
	if _, ok := p.storage.Data[unsubscribeData.EventName]; !ok {
		errorText := fmt.Sprintf("No subscribers for event %s", unsubscribeData.EventName)
		return errors.New(errorText)
	}
	for pos, element := range p.storage.Data[unsubscribeData.EventName] {
		if element == unsubscribeData.Url {
			p.storage.Data[unsubscribeData.EventName][pos] = p.storage.Data[unsubscribeData.EventName][len(p.storage.Data[unsubscribeData.EventName])-1]
			p.storage.Data[unsubscribeData.EventName] = p.storage.Data[unsubscribeData.EventName][:len(p.storage.Data[unsubscribeData.EventName])-1]
			if len(p.storage.Data[unsubscribeData.EventName]) == 0 {
				delete(p.storage.Data, unsubscribeData.EventName)
			}
			p.storage.Mutex.Unlock()
			return nil
		}
	}

	p.storage.Mutex.Unlock()

	return errors.New("No subscribers found")
}

func (p PluginWeb) ProcessEvent(eventData event.Event) {
	httpClient := &http.Client{}
	for _, subscribeUrl := range p.storage.Data[eventData.Name] {
		req, err := http.NewRequest("POST", subscribeUrl, bytes.NewBuffer(eventData.Payload))
		if err != nil {
			log.Error(err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Println("Status of sending event", eventData.Name, "is", resp.Status)
	}
}
