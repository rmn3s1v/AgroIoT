package validator

import (
	"errors"
	"fmt"
	"time"

	"ingestion-service/internal/model"
)

type MetricRule struct {
	Min float64
	Max float64
}

var metricRules = map[string]map[string]MetricRule{
	"soil_sensor": {
		"soil_moisture":   {Min: 0, Max: 100},
		"soil_temperature": {Min: -30, Max: 60},
		"soil_ph":         {Min: 0, Max: 14},
	},
	"climate_sensor": {
		"temperature":  {Min: -50, Max: 60},
		"air_humidity": {Min: 0, Max: 100},
		"lux":          {Min: 0, Max: 100000},
		"co2_level":    {Min: 0, Max: 5000},
	},
}

var requiredMetrics = map[string][]string{
	"soil_sensor": {
		"soil_moisture",
	},
	"climate_sensor": {
		"temperature",
		"air_humidity",
	},
}

func ValidateTelemetry(t model.Telemetry) error {
	if t.DeviceID == "" {
		return errors.New("device_id is required")
	}

	if t.DeviceType == "" {
		return errors.New("device_type is required")
	}

	rules, ok := metricRules[t.DeviceType]
	if !ok {
		return fmt.Errorf("unknown device_type: %s", t.DeviceType)
	}

	if t.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	if t.Timestamp.After(time.Now().Add(1 * time.Minute)) {
		return errors.New("timestamp is in the future")
	}

	if len(t.Metrics) == 0 {
		return errors.New("metrics cannot be empty")
	}

	for _, required := range requiredMetrics[t.DeviceType] {
		if _, ok := t.Metrics[required]; !ok {
			return fmt.Errorf("missing required metric: %s", required)
		}
	}

	for key, value := range t.Metrics {
		rule, ok := rules[key]
		if !ok {
			return fmt.Errorf("invalid metric '%s' for device_type '%s'", key, t.DeviceType)
		}

		num, err := toFloat(value)
		if err != nil {
			return fmt.Errorf("metric '%s' must be numeric", key)
		}

		if num < rule.Min || num > rule.Max {
			return fmt.Errorf("metric '%s' out of range [%f, %f]", key, rule.Min, rule.Max)
		}
	}

	return nil
}

func toFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	default:
		return 0, errors.New("not a number")
	}
}

func NormalizeMetrics(metrics map[string]interface{}) map[string]float64 {
	result := make(map[string]float64)

	for k, v := range metrics {
		if val, err := toFloat(v); err == nil {
			result[k] = val
		}
	}

	return result
}
