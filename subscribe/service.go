package subscribe

import (
	"errors"
	log "github.com/Sirupsen/logrus"

	"github.com/qa-dev/universe/storage"
)

type SubscribeService struct {
	storage *storage.Storage
}

func NewSubscribeService(storage *storage.Storage) *SubscribeService {
	return &SubscribeService{storage}
}

func (s *SubscribeService) ProcessSubscribe(subscribe Subscribe) error {
	if len(subscribe.EventName) == 0 {
		log.Println("Got blank subscribe event name")
		return errors.New("BLANK EVENT NAME")
	}
	if len(subscribe.WebHookPath) == 0 {
		log.Println("Got blank subscribe webhook path")
		return errors.New("BLANK WEBHOOK PATH")
	}
	s.storage.Mutex.Lock()
	if _, ok := s.storage.Data[subscribe.EventName]; !ok {
		s.storage.Data[subscribe.EventName] = make([]string, 0)
	}
	s.storage.Data[subscribe.EventName] = append(s.storage.Data[subscribe.EventName], subscribe.WebHookPath)
	s.storage.Mutex.Unlock()

	log.Println("Subscribed on", subscribe.EventName, "to", subscribe.WebHookPath)

	return nil
}
