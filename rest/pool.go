package rest

import (
	"net/http"
	"sync"
)

type TransportCacheInterfce interface {
	Get(config *TransportConfig) *http.Transport
}

var TransportCache = newTransportPool()

type transportCacheKey string

func (t transportCacheKey) String() string {
	return ""
}

type transportPool struct {
	sync.RWMutex
	pool map[transportCacheKey]*http.Transport
}

func (t *transportPool) Get(config *TransportConfig) *http.Transport {
	return nil
}

var _ TransportCacheInterfce = &transportPool{}

func newTransportPool() *transportPool {
	return &transportPool{
		pool: make(map[transportCacheKey]*http.Transport),
	}
}
