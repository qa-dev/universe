package web

import "sync"

type Storage struct {
	Data  map[string][]string
	Mutex *sync.Mutex
}

func NewStorage() *Storage {
	dt := make(map[string][]string)
	mtx := &sync.Mutex{}
	return &Storage{Data: dt, Mutex: mtx}
}
