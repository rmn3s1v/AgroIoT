package handler

import (
	"api-service/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.TelemetryService
}

func NewHandler(s *service.TelemetryService) *Handler {
	return &Handler{Service: s}
}

// GetTelemetry returns telemetry records for a device.
// @Summary Get device telemetry
// @Description Returns telemetry records for a device in the requested time range.
// @Tags telemetry
// @Produce json
// @Param device_id path string true "Device ID"
// @Param from query string false "Start time in RFC3339 format" example(2026-05-10T00:00:00Z)
// @Param to query string false "End time in RFC3339 format" example(2026-05-10T23:59:59Z)
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /devices/{device_id}/telemetry [get]
func (h *Handler) GetTelemetry(c *gin.Context) {
	deviceID := c.Param("device_id")

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, _ := time.Parse(time.RFC3339, fromStr)
	to, _ := time.Parse(time.RFC3339, toStr)

	data, err := h.Service.GetByDevice(deviceID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetAggregates returns aggregate values for a device metric.
// @Summary Get metric aggregates
// @Description Returns avg, min, and max values for a device metric in the requested time range.
// @Tags telemetry
// @Produce json
// @Param device_id path string true "Device ID"
// @Param metric query string true "Metric name" example(temperature)
// @Param from query string false "Start time in RFC3339 format" example(2026-05-10T00:00:00Z)
// @Param to query string false "End time in RFC3339 format" example(2026-05-10T23:59:59Z)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /devices/{device_id}/aggregates [get]
func (h *Handler) GetAggregates(c *gin.Context) {
	deviceID := c.Param("device_id")
	metric := c.Query("metric")

	if metric == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric is required"})
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, _ := time.Parse(time.RFC3339, fromStr)
	to, _ := time.Parse(time.RFC3339, toStr)

	result, err := h.Service.GetAggregates(deviceID, metric, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no data"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTimeSeriesAggregates returns time-series aggregate values for a device metric.
// @Summary Get time-series metric aggregates
// @Description Returns avg, min, and max values grouped by hour or day for a device metric.
// @Tags telemetry
// @Produce json
// @Param device_id path string true "Device ID"
// @Param metric query string true "Metric name" example(temperature)
// @Param interval query string false "Aggregation interval" Enums(hour, day) default(hour)
// @Param from query string false "Start time in RFC3339 format" example(2026-05-10T00:00:00Z)
// @Param to query string false "End time in RFC3339 format" example(2026-05-10T23:59:59Z)
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /devices/{device_id}/aggregates/timeseries [get]
func (h *Handler) GetTimeSeriesAggregates(c *gin.Context) {
	deviceID := c.Param("device_id")
	metric := c.Query("metric")
	interval := c.Query("interval")

	if metric == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric is required"})
		return
	}

	if interval == "" {
		interval = "hour"
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, _ := time.Parse(time.RFC3339, fromStr)
	to, _ := time.Parse(time.RFC3339, toStr)

	result, err := h.Service.GetTimeSeriesAggregates(deviceID, metric, interval, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
