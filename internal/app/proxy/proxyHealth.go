package proxy

import (
	"net"
	"net/url"
	"sync"
	"time"
)

type ProxyHealth struct {
	origin *url.URL

	mutex         sync.Mutex
	period        time.Duration
	cancel        chan struct{}
	isAvailable   bool
	timeoutPeriod time.Duration
}

func NewProxyHealth(origin *url.URL, period, timeoutPeriod time.Duration) *ProxyHealth {
	h := &ProxyHealth{
		origin:        origin,
		period:        period,
		cancel:        make(chan struct{}),
		isAvailable:   checkHealth(origin, period),
		timeoutPeriod: timeoutPeriod,
	}
	h.Start()
	return h
}

func (h *ProxyHealth) Start() {
	go func() {
		t := time.NewTicker(h.period)
		for {
			select {
			case <-t.C:
				h.isAvailable = checkHealth(h.origin, h.period)
			case <-h.cancel:
				t.Stop()
				return
			}
		}
	}()
}

func checkHealth(origin *url.URL, period time.Duration) bool {
	conn, err := net.DialTimeout("tcp", origin.Host, period)
	if err != nil {
		//log.Println(err)
		return false
	}
	_ = conn.Close()
	return true
}

func (h *ProxyHealth) IsAvailable() bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.isAvailable
}
