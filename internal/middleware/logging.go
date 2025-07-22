package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLoggerFormatter
	Formatter gin.LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// responseBodyWriter is a custom writer to capture response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger returns a gin.HandlerFunc (middleware) that logs requests using logrus.
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig returns a gin.HandlerFunc using configs.
func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: formatter,
		Output:    out,
		SkipPaths: notlogged,
	})
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	// Create structured log entry
	entry := logrus.WithFields(logrus.Fields{
		"timestamp":   param.TimeStamp.Format(time.RFC3339),
		"client_ip":   param.ClientIP,
		"method":      param.Method,
		"path":        param.Path,
		"protocol":    param.Request.Proto,
		"status_code": param.StatusCode,
		"latency":     param.Latency.String(),
		"user_agent":  param.Request.UserAgent(),
		"data_length": param.BodySize,
	})

	// Add error message if present
	if param.ErrorMessage != "" {
		entry = entry.WithField("error", param.ErrorMessage)
	}

	// Add request ID if present
	if requestID := param.Request.Header.Get("X-Request-ID"); requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}

	// Log based on status code
	switch {
	case param.StatusCode >= 500:
		entry.Error("HTTP Request")
	case param.StatusCode >= 400:
		entry.Warn("HTTP Request")
	default:
		entry.Info("HTTP Request")
	}

	// Return formatted string for gin's default output
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

//// StructuredLogger returns a structured logger middleware
//func StructuredLogger() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Start timer
//		start := time.Now()
//		path := c.Request.URL.Path
//		raw := c.Request.URL.RawQuery
//
//		// Capture request body for logging (optional, be careful with large payloads)
//		var requestBody []byte
//		if c.Request.Body != nil && c.Request.ContentLength > 0 && c.Request.ContentLength < 1024 { // Only log small payloads
//			requestBody, _ = io.ReadAll(c.Request.Body)
//			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
//		}
//
//		// Capture response body
//		responseBody := &bytes.Buffer{}
//		writer := &responseBodyWriter{
//			ResponseWriter: c.Writer,
//			body:           responseBody,
//		}
//		c.Writer = writer
//
//		// Process request
//		c.Next()
//
//		// Calculate latency
//		latency := time.Since(start)
//
//		// Get client IP
//		clientIP := c.ClientIP()
//
//		// Get method
//		method := c.Request.Method
//
//		// Get status code
//		statusCode := c.Writer.Status()
//
//		// Build full path
//		if raw != "" {
//			path = path + "?" + raw
//		}
//
//		// Create structured log entry
//		fields := logrus.Fields{
//			"timestamp":     start.Format(time.RFC3339),
//			"client_ip":     clientIP,
//			"method":        method,
//			"path":          path,
//			"protocol":      c.Request.Proto,
//			"status_code":   statusCode,
//			"latency_ms":    float64(latency.Nanoseconds()) / 1e6,
//			"user_agent":    c.Request.UserAgent(),
//			"response_size": c.Writer.Size(),
//		}
//
//		// Add request ID if present
//		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
//			fields["request_id"] = requestID
//		}
//
//		// Add user info if present
//		if userID, exists := c.Get("user_id"); exists {
//			fields["user_id"] = userID
//		}
//
//		// Add request body for certain content types (be selective)
//		if len(requestBody) > 0 && c.GetHeader("Content-Type") == "application/json" {
//			fields["request_body"] = string(requestBody)
//		}
//
//		// Add response body for errors (be selective)
//		if statusCode >= 400 && responseBody.Len() < 500 {
//			fields["response_body"] = responseBody.String()
//		}
//
//		// Get any errors that occurred during request processing
//		if len(c.Errors) > 0 {
//			fields["errors"] = c.Errors.String()
//		}
//
//		// Log based on status code
//		entry := logger.WithFields(fields)
//		switch {
//		case statusCode >= 500:
//			entry.Error("HTTP Request - Server Error")
//		case statusCode >= 400:
//			entry.Warn("HTTP Request - Client Error")
//		case statusCode >= 300:
//			entry.Info("HTTP Request - Redirect")
//		default:
//			entry.Info("HTTP Request - Success")
//		}
//	}
//}
//
//// RequestIDMiddleware generates and adds request ID to context
//func RequestIDMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		requestID := c.GetHeader("X-Request-ID")
//		if requestID == "" {
//			requestID = generateRequestID()
//		}
//
//		c.Header("X-Request-ID", requestID)
//		c.Set("request_id", requestID)
//		c.Next()
//	}
//}
//
//// AccessLogger logs only access information (lighter than StructuredLogger)
//func AccessLogger() gin.HandlerFunc {
//	return gin.LoggerWithConfig(gin.LoggerConfig{
//		Formatter: func(param gin.LogFormatterParams) string {
//			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
//				param.ClientIP,
//				param.TimeStamp.Format(time.RFC1123),
//				param.Method,
//				param.Path,
//				param.Request.Proto,
//				param.StatusCode,
//				param.Latency,
//				param.Request.UserAgent(),
//				param.ErrorMessage,
//			)
//		},
//		Output:    logger.GetOutput(),
//		SkipPaths: []string{"/health", "/metrics"},
//	})
//}

//// generateRequestID generates a unique request ID
//func generateRequestID() string {
//	// Simple implementation - in production, consider using UUID
//	return fmt.Sprintf("%d", time.Now().UnixNano())
//}
