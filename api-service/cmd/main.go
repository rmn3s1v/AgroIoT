package main

import (
	"api-service/internal/handler"
	"api-service/internal/mongodb"
	"api-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	mongo := mongodb.NewClient("mongodb://localhost:27017")

	service := service.NewTelemetryService(mongo.Collection)
	handler := handler.NewHandler(service)

	r := gin.Default()

	r.GET("/devices/:device_id/telemetry", handler.GetTelemetry)
	r.GET("/devices/:device_id/aggregates", handler.GetAggregates)
	r.GET("/devices/:device_id/aggregates/timeseries", handler.GetTimeSeriesAggregates)

	log.Println("API service running on :8080")
	r.Run(":8080")
}
