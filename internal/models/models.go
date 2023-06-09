package models

import "time"

type WebsiteCheck struct {
	Access       bool
	ResponseTime time.Duration
}

type PingResult struct {
	Data map[string]WebsiteCheck
	Min  string
	Max  string
}

type Stats struct {
	Endpoint string
	Counter  uint64
}
