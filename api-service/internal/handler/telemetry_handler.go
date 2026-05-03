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
