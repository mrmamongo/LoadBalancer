package balancer

import (
	"fmt"
	"github.com/mrmamongo/LoadBalancer/internal/app/proxy"
	"net/http"
	"sort"
)

type Balancer struct {
	proxies []*proxy.Proxy
}

func NewBalancer(proxies ...*proxy.Proxy) *Balancer {
	return &Balancer{
		proxies: proxies,
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pr, err := b.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("No response from server"))
		//if err != nil {
		//log.Fatal(err)
		//}
	}
	pr.ServeHTTP(w, r)
}

func (b *Balancer) Next() (*proxy.Proxy, error) {
	if len(b.proxies) == 0 {
		return nil, fmt.Errorf("no proxies available")
	}
	proxies := append(b.proxies[:0:0], b.proxies...)
	sort.SliceStable(proxies, func(i, j int) bool {
		return proxies[i].Connections() < proxies[j].Connections()
	})
	for _, pr := range proxies {
		if pr.Available() {
			return pr, nil
		}
	}
	return nil, fmt.Errorf("no proxies available")
}
