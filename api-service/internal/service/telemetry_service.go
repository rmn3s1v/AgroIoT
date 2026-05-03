package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TelemetryService struct {
	Collection *mongo.Collection
}

func NewTelemetryService(c *mongo.Collection) *TelemetryService {
	return &TelemetryService{Collection: c}
}

func (s *TelemetryService) GetByDevice(deviceID string, from, to time.Time) ([]bson.M, error) {
	filter := bson.M{
		"device_id": deviceID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	cursor, err := s.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *TelemetryService) GetAggregates(deviceID, metric string, from, to time.Time) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"device_id": deviceID,
				"timestamp": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"avg": bson.M{"$avg": "$metrics." + metric},
				"min": bson.M{"$min": "$metrics." + metric},
				"max": bson.M{"$max": "$metrics." + metric},
			},
		},
	}

	cursor, err := s.Collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

func (s *TelemetryService) GetTimeSeriesAggregates(
	deviceID, metric, interval string,
	from, to time.Time,
) ([]bson.M, error) {

	var dateFormat string

	switch interval {
	case "hour":
		dateFormat = "%Y-%m-%d %H:00"
	case "day":
		dateFormat = "%Y-%m-%d"
	default:
		return nil, fmt.Errorf("invalid interval: %s", interval)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"device_id": deviceID,
				"timestamp": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"time": bson.M{
						"$dateToString": bson.M{
							"format": dateFormat,
							"date":   "$timestamp",
						},
					},
				},
				"avg": bson.M{"$avg": "$metrics." + metric},
				"min": bson.M{"$min": "$metrics." + metric},
				"max": bson.M{"$max": "$metrics." + metric},
			},
		},
		{
			"$sort": bson.M{"_id.time": 1},
		},
	}

	cursor, err := s.Collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
