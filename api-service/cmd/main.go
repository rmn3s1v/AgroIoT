package main

import (
	_ "api-service/docs"
	"api-service/internal/handler"
	"api-service/internal/mongodb"
	"api-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AgroIoT API Service
// @version 1.0
// @description API for reading device telemetry and aggregate metrics.
// @host localhost:8080
// @BasePath /
func main() {
	mongo := mongodb.NewClient("mongodb://mongodb:27017")

	service := service.NewTelemetryService(mongo.Collection)
	handler := handler.NewHandler(service)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/devices/:device_id/telemetry", handler.GetTelemetry)
	r.GET("/devices/:device_id/aggregates", handler.GetAggregates)
	r.GET("/devices/:device_id/aggregates/timeseries", handler.GetTimeSeriesAggregates)

	log.Println("API service running on :8080")
	r.Run(":8080")
}
