package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mrmamongo/LoadBalancer/internal/app/balancer"
	"github.com/mrmamongo/LoadBalancer/internal/app/config"
	"github.com/mrmamongo/LoadBalancer/internal/app/proxy"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	file, _ := os.Open(*configPath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	configuration := config.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln(err)
	}

	var proxies []*proxy.Proxy
	var isAnyProxyAvailable bool
	for _, pr := range configuration.Proxies {
		origin, e := url.Parse(pr.Host)
		if e != nil {
			log.Fatalln(e)
		}
		p := proxy.NewProxy(pr.Index, origin, pr.Period*time.Second, pr.Timeout*time.Second)
		if p.Available() {
			log.Println("Proxy " + pr.Index + " started")
			isAnyProxyAvailable = true
		} else {
			log.Println("Proxy " + pr.Index + " is not available")
		}
		proxies = append(proxies, p)
	}
	if !isAnyProxyAvailable {
		log.Fatalln("No available proxies")
	}

	b := balancer.NewBalancer(proxies...)
	log.Println("Server started, config:", configuration)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", configuration.Host, configuration.Port), b)
	if err != nil {
		log.Fatalln(err)
	}

}
