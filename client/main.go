package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type ConnectionWatcher struct {
	n int64
}

func NewConnectionWatcher() *ConnectionWatcher {
	return &ConnectionWatcher{}
}

func (cw *ConnectionWatcher) OnStateChange(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		cw.Add(1)
	case http.StateHijacked, http.StateClosed:
		cw.Add(-1)
	}
}

func (cw *ConnectionWatcher) Connections() int {
	return int(atomic.LoadInt64(&cw.n))
}

func (cw *ConnectionWatcher) Add(c int64) {
	atomic.AddInt64(&cw.n, c)
}

func startLogger(index string, quitChan chan struct{}, watcher *ConnectionWatcher, server *http.Server) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Server: "+index+" connections: ", watcher.Connections())
			case <-quitChan:
				log.Println("Shutting down server")
				ticker.Stop()
				err := server.Shutdown(nil)
				if err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	}()
}

type HelloHandler struct {
	delay int
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Duration(h.delay) * time.Second)
	_, err := fmt.Fprintf(w, "Hello!")
	if err != nil {
		return
	}
}

type QuitHandler struct {
	quitHandler chan<- struct{}
}

func (qh *QuitHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Quitting!")
	if err != nil {
		return
	}
	qh.quitHandler <- struct{}{}
}

type HealthCheckHandler struct{}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		return
	}
}

func main() {
	index := flag.String("index", "test", "Server index")
	host := flag.String("host", "localhost", "Host to listen on")
	port := flag.Int("port", 25565, "Port to listen on")
	delay := flag.Int("delay", 1, "Delay to answer to requests")
	flag.Parse()

	watcher := NewConnectionWatcher()
	quitChan := make(chan struct{})

	mux := http.NewServeMux()
	mux.Handle("/hc", &HealthCheckHandler{})
	mux.Handle("/hello", &HelloHandler{*delay})
	mux.Handle("/quit", &QuitHandler{quitChan})

	server := http.Server{
		ConnState: watcher.OnStateChange,
		Addr:      *host + ":" + strconv.Itoa(*port),
		Handler:   mux,
	}
	go startLogger(*index, quitChan, watcher, &server)
	log.Println("Starting server on " + server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
