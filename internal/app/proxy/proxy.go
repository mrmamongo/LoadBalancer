package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Proxy struct {
	health *ProxyHealth

	load  int32
	proxy *httputil.ReverseProxy
	index string
}

func NewProxy(index string, origin *url.URL, period, timeoutPeriod time.Duration) *Proxy {
	reverseProxy := httputil.NewSingleHostReverseProxy(origin)
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
	}
	return &Proxy{
		index:  index,
		proxy:  reverseProxy,
		health: NewProxyHealth(origin, period, timeoutPeriod),
	}
}

func (p *Proxy) Index() string {
	return p.index
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&p.load, 1)
	defer atomic.AddInt32(&p.load, -1)
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) Connections() int32 {
	return atomic.LoadInt32(&p.load)
}

func (p *Proxy) Available() bool {
	return p.health.IsAvailable()
}
