package config

import (
	"sync"
)

type defaultConfigProvider struct {
	mu   sync.RWMutex
	data map[string]any
}

var DefaultProvider Provider = &defaultConfigProvider{
	data: make(map[string]any),
}
