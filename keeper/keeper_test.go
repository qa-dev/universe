package keeper

import (
	"encoding/json"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
)

type TestSubscribeInfoType struct {
	Hello string
}

type FakePlugin struct {
	subscribers []string
}

func (m FakePlugin) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		Name:                "Name",
		Tag:                 "test_fake_plugin",
		Version:             1,
		SubscribersStorable: true,
	}
}

func (FakePlugin) Subscribe(input []byte) error {
	return nil
}

func (FakePlugin) Unsubscribe(input []byte) error {
	return nil
}

func (FakePlugin) ProcessEvent(eventData event.Event) {}

func (p *FakePlugin) LoadSubscriber(data []byte) error {
	var subscribeData TestSubscribeInfoType
	err := json.Unmarshal(data, &subscribeData)
	if err != nil {
		return err
	}
	p.subscribers = append(p.subscribers, subscribeData.Hello)
	return nil
}

func (FakePlugin) Loaded() {}
