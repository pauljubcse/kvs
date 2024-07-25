package kvs

import "sync"

type Domain struct {
	stringStore   map[string]string
	skipListStore map[string]*SkipList
	mu            sync.RWMutex
}

func NewDomain() *Domain {
	return &Domain{
		stringStore:   make(map[string]string),
		skipListStore: make(map[string]*SkipList),
	}
}