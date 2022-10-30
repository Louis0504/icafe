package redis

import (
	"fmt"
	"sync"
)

var RWRedisManager = newGroupManager()

type GroupManager struct {
	mu     sync.RWMutex
	rwredis map[string]*RWRedis
}

func newGroupManager() *GroupManager {
	return &GroupManager{
		rwredis: make(map[string]*RWRedis),
	}
}

func (gm *GroupManager) Add(name string, g *RWRedis) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if _, ok := gm.rwredis[name]; ok {
		return fmt.Errorf("redis group alread exists")
	}
	gm.rwredis[name] = g
	return nil
}

func (gm *GroupManager) Get(name string) *RWRedis {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	return gm.rwredis[name]
}

func Get(name string) *RWRedis {
	return RWRedisManager.Get(name)
}
