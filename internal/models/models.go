package models

import "time"

type WebsiteCheck struct {
	Access       bool
	ResponseTime time.Duration
}
