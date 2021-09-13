package id_generator

import (
	"sync"
	"sync/atomic"
)

type IdGenerator interface {
	GenerateId() int32
	GetLastId() int32
}

type idGenerator struct {
	curId int32
}

var (
	instance *idGenerator
	lock     = &sync.Mutex{}
)

func GetInstance() IdGenerator {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = &idGenerator{}
	}
	return instance
}

func (s *idGenerator) GenerateId() int32 {
	return atomic.AddInt32(&s.curId, 1)
}

func (s *idGenerator) GetLastId() int32 {
	return atomic.LoadInt32(&s.curId)
}
