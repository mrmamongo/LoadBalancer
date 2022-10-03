package config

import "time"

type ConfigProxy struct {
	Index   string        `json:"index"`
	Host    string        `json:"host"`
	Period  time.Duration `json:"period"`
	Timeout time.Duration `json:"timeout"`
}

type Configuration struct {
	Port    int           `json:"port"`
	Host    string        `json:"host"`
	Proxies []ConfigProxy `json:"proxies"`
}
