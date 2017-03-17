package subscribe

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/plugins"
)

type SubscribeService struct {
	pluginStorage *plugins.PluginStorage
}

func NewSubscribeService(storage *plugins.PluginStorage) *SubscribeService {
	return &SubscribeService{storage}
}

func (s *SubscribeService) ProcessSubscribe(pluginName string, input []byte) error {
	if len(pluginName) == 0 {
		log.Println("Got blank plugin name")
		return errors.New("BLANK PLUGIN NAME")
	}
	err := s.pluginStorage.ProcessSubscribe(pluginName, input)

	return err
}

func (s *SubscribeService) ProcessUnsubscribe(pluginName string, input []byte) error {
	if len(pluginName) == 0 {
		log.Println("Got blank plugin name")
		return errors.New("BLANK PLUGIN NAME")
	}
	err := s.pluginStorage.ProcessUnsubscribe(pluginName, input)

	return err
}
