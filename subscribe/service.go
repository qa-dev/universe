package subscribe

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/plugins"
)

type SubscribeService struct{}

func NewSubscribeService() *SubscribeService {
	return &SubscribeService{}
}

func (s *SubscribeService) ProcessSubscribe(pluginName string, input []byte) error {
	if len(pluginName) == 0 {
		log.Println("Got blank plugin name")
		return errors.New("BLANK PLUGIN NAME")
	}
	err := plugins.Obs.ProcessSubscribe(pluginName, input)

	return err
}
