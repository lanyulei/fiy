package cache

import (
	"fiy/pkg/core/cache"
)

var MemoryAdapter Adapter

func InitMemory() error {
	MemoryAdapter = &cache.Memory{
		PoolNum: 100,
	}
	err := MemoryAdapter.Connect()
	return err
}
