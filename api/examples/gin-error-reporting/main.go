package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// TracClient represents a simple client to send errors to trac API
type TracClient struct {
	BaseURL   string
	PublicKey string
	ProjectID string
	Client    *http.Client
}

// Event represents an error event to send to trac
type Event struct {
	EventID   string                 `json:"event_id"`
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Platform  string                 `json:"platform"`
	SDK       map[string]interface{} `json:"sdk"`
	Exception Exception              `json:"exception"`
	Tags      map[string]string      `json:"tags"`
	User      map[string]interface{} `json:"user"`
	Request   Request                `json:"request"`
}

type Exception struct {
	Type       string       `json:"type"`
	Value      string       `json:"value"`
	Stacktrace []StackFrame `json:"stacktrace"`
}

type StackFrame struct {
	Filename     string `json:"filename"`
	Function     string `json:"function"`
	LineNumber   int    `json:"lineno"`
	ColumnNumber int    `json:"colno"`
}

type Request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type Envelope struct {
	EventID string         `json:"event_id"`
	Items   []EnvelopeItem `json:"items"`
}

type EnvelopeItem struct {
	Header  map[string]string `json:"headers"`
	Type    string            `json:"type"`
	Payload json.RawMessage   `json:"payload"`
}

// NewTracClient creates a new trac client
func NewTracClient(dsn string) *TracClient {
	// Parse DSN: http://publickey@host:port/projectid
	var publicKey, host, projectID string
	fmt.Sscanf(dsn, "http://%s@%s/%s", &publicKey, &host, &projectID)

	return &TracClient{
		BaseURL:   fmt.Sprintf("http://%s", host),
		PublicKey: publicKey,
		ProjectID: projectID,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// CaptureException sends an exception to trac
func (t *TracClient) CaptureException(err error, c *gin.Context) string {
	eventID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Create stack trace
	stack := make([]StackFrame, 0)
	for _, pc := range callStack() {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		stack = append(stack, StackFrame{
			Filename:     file,
			Function:     fn.Name(),
			LineNumber:   line,
			ColumnNumber: 0,
		})
	}

	// Prepare headers
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	event := Event{
		EventID:   eventID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "error",
		Message:   err.Error(),
		Platform:  "go",
		SDK: map[string]interface{}{
			"name":    "trac-go",
			"version": "1.0.0",
		},
		Exception: Exception{
			Type:       "panic",
			Value:      err.Error(),
			Stacktrace: stack,
		},
		Tags: map[string]string{
			"endpoint": c.Request.URL.Path,
			"method":   c.Request.Method,
		},
		User: map[string]interface{}{
			"id":       "123",
			"username": "testuser",
		},
		Request: Request{
			URL:     c.Request.URL.String(),
			Method:  c.Request.Method,
			Headers: headers,
		},
	}

	// Create envelope
	eventPayload, _ := json.Marshal(event)
	envelope := Envelope{
		EventID: eventID,
		Items: []EnvelopeItem{
			{
				Header: map[string]string{
					"type": "event",
				},
				Type:    "event",
				Payload: json.RawMessage(eventPayload),
			},
		},
	}

	// Send to trac
	go func() {
		envelopeData, _ := json.Marshal(envelope)
		url := fmt.Sprintf("%s/api/%s/envelope/", t.BaseURL, t.ProjectID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(envelopeData))
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Sentry-Auth", fmt.Sprintf("Sentry sentry_version=7, sentry_key=%s", t.PublicKey))

		resp, err := t.Client.Do(req)
		if err != nil {
			log.Printf("Failed to send error to trac: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("Trac returned status %d", resp.StatusCode)
		}
	}()

	return eventID
}

func callStack() []uintptr {
	const size = 32
	var pcs [size]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

func main() {
	// Initialize trac client with DSN from environment or use default
	dsn := os.Getenv("TRAC_DSN")
	if dsn == "" {
		dsn = "http://44da05e841a422acd8b414dbb452e88a@localhost:8025/2"
	}

	tracClient := NewTracClient(dsn)

	// Create Gin router
	r := gin.Default()

	// Add a middleware to recover from panics and capture errors
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Convert panic to error
				var panicErr error
				switch x := err.(type) {
				case string:
					panicErr = fmt.Errorf(x)
				case error:
					panicErr = x
				default:
					panicErr = fmt.Errorf("unknown panic: %v", x)
				}

				// Capture the error with trac
				eventID := tracClient.CaptureException(panicErr, c)

				// Log the stack trace
				log.Printf("Panic recovered: %v\n%s", panicErr, debug.Stack())

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Internal Server Error",
					"event_id": eventID,
					"message":  "Error has been reported to trac",
				})
				c.Abort()
			}
		}()

		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Capture the error with trac
			eventID := tracClient.CaptureException(err.Err, c)

			// Return error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Internal Server Error",
				"event_id": eventID,
				"message":  "Error has been reported to trac",
			})
		}
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"version":   "1.0.0",
		})
	})

	// Test endpoints
	r.GET("/test/panic", func(c *gin.Context) {
		// This will cause a panic
		panic("This is a test panic!")
	})

	r.GET("/test/error", func(c *gin.Context) {
		// This will cause an error
		c.Error(fmt.Errorf("This is a test error"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong"})
	})

	r.GET("/test/divide", func(c *gin.Context) {
		// This will cause a division by zero error
		a := 10
		b := 0
		result := a / b
		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	r.POST("/test/process", func(c *gin.Context) {
		var req struct {
			Name string `json:"name" binding:"required"`
			Age  int    `json:"age" binding:"required,min=0,max=150"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		// Simulate processing
		if req.Name == "error" {
			c.Error(fmt.Errorf("processing failed for name: %s", req.Name))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Processed %s, age %d", req.Name, req.Age),
			"success": true,
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Gin server starting on port %s\n", port)
	fmt.Printf("📊 Reporting errors to trac: %s\n", dsn)
	fmt.Printf("🌍 Environment: %s\n", os.Getenv("GIN_ENV"))
	fmt.Printf("🔧 Go Version: %s\n", runtime.Version())

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
