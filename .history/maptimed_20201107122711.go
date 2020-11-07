package maptimed

import (
	"errors"
	"sync"
	"time"
)

// MapTimed is a structure with key (string) -> value (any type) which erases all keys that haven't been accessed after a specified timeout
type MapTimed struct {
	m       map[string]*mapTimedValue
	timeout time.Duration
	timer   *time.Timer
	lock    sync.Mutex
}

type mapTimedValue struct {
	v       interface{}
	laccess time.Time
}

// NewMapTimed returns a MapTimed structure that will delete all keys that haven't been accessed in more than timeout duration
func NewMapTimed(timeout time.Duration) (*MapTimed, error) {
	if timeout < 1 {
		return nil, errors.New("Please provide a timeout value of at least 1s")
	}

	mt := new(MapTimed)
	mt.m = make(map[string]*mapTimedValue)
	mt.timeout = timeout
	mt.timer = time.NewTimer(mt.timeout * time.Second)
	go mt.clear()
	return mt, nil
}

// Get a value from a MapTimed
func (mt *MapTimed) Get(key string) (val interface{}) {
	mt.lock.Lock()
	defer mt.lock.Unlock()

	_, exists := mt.m[key]
	if !exists {
		return nil
	}
	mt.m[key].laccess = time.Now()

	return mt.m[key].v
}

// Set a value in a MapTimed
func (mt *MapTimed) Set(key string, val interface{}) {
	mt.lock.Lock()
	defer mt.lock.Unlock()

	_, exists := mt.m[key]
	if !exists {
		mt.m[key] = new(mapTimedValue)
	}
	mt.m[key].v = val
	mt.m[key].laccess = time.Now()
}

func (mt *MapTimed) clear() {
	<-mt.timer.C
	for key, mtv := range mt.m {
		if time.Since(mtv.laccess) > mt.timeout {
			mt.lock.Lock()
			delete(mt.m, key)
			mt.lock.Unlock()
		}
	}
	mt.timer.Reset(mt.timeout)
	mt.clear()
}
