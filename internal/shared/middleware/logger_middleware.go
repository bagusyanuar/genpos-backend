package middleware

import (
	"time"

	"github.com/bagusyanuar/genpos-backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Logger(conf *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Skip logging for health checks or high-volume monitoring endpoints
		path := c.Path()
		if path == "/health" || path == "/metrics" {
			return c.Next()
		}

		start := time.Now()

		// Go to next middleware/handler
		err := c.Next()

		latency := time.Since(start)
		requestID := c.GetRespHeader(fiber.HeaderXRequestID)
		if requestID == "" {
			requestID = c.Get(fiber.HeaderXRequestID)
		}

		// Prepare fields for logging
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", path),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("latency", latency.String()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
		}

		// Log based on status code
		msg := "Request Completed"
		if err != nil {
			fields = append(fields, zap.Error(err))
			config.Log.Error("Request Failed", fields...)
		} else if c.Response().StatusCode() >= 400 {
			config.Log.Warn(msg, fields...)
		} else {
			config.Log.Info(msg, fields...)
		}

		return err
	}
}
