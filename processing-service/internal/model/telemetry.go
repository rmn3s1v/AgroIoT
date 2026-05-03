package model

import "time"

type Telemetry struct {
	DeviceID   string             `bson:"device_id" json:"device_id"`
	DeviceType string             `bson:"device_type" json:"device_type"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	Metrics    map[string]float64 `bson:"metrics" json:"metrics"`
}
