package src

import (
	"sync"
)

type ConnectionManager struct {
	Pools map[string]Pool
}

var once sync.Once
var instance *ConnectionManager

func GetManager() *ConnectionManager {
	once.Do(func() {
		instance = &ConnectionManager{Pools: make(map[string]Pool, 10)}
	})
	return instance
}
