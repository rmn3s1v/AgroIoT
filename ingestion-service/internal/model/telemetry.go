package model

import "time"

type Telemetry struct {
	DeviceID   string                 `json:"device_id"`
	DeviceType string                 `json:"device_type"`
	Timestamp  time.Time              `json:"timestamp"`
	Metrics    map[string]interface{} `json:"metrics"`
}
