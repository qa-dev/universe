package web

import (
	"errors"
	"fmt"
	"sync"
)

type Storage struct {
	Data  map[string][]string
	mutex *sync.Mutex
}

func NewStorage() *Storage {
	dt := make(map[string][]string)
	mtx := &sync.Mutex{}
	return &Storage{Data: dt, mutex: mtx}
}

func (s *Storage) AddSubscriber(data SubscribeData) {
	s.mutex.Lock()
	if _, ok := s.Data[data.EventName]; !ok {
		s.Data[data.EventName] = make([]string, 0)
	}
	s.Data[data.EventName] = append(s.Data[data.EventName], data.Url)
	s.mutex.Unlock()
}

func (s *Storage) RemoveSubscriber(data UnsubscribeData) error {
	s.mutex.Lock()
	if _, ok := s.Data[data.EventName]; !ok {
		errorText := fmt.Sprintf("No subscribers for event %s", data.EventName)
		return errors.New(errorText)
	}
	for pos, element := range s.Data[data.EventName] {
		if element == data.Url {
			s.Data[data.EventName][pos] = s.Data[data.EventName][len(s.Data[data.EventName])-1]
			s.Data[data.EventName] = s.Data[data.EventName][:len(s.Data[data.EventName])-1]
			if len(s.Data[data.EventName]) == 0 {
				delete(s.Data, data.EventName)
			}
			s.mutex.Unlock()
			return nil
		}
	}
	s.mutex.Unlock()
	return errors.New("No subscribers found")
}
