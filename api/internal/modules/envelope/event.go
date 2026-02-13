package envelope

import (
	"encoding/json"
	"time"
)

// Event represents a Sentry-compatible error event
// This is the main data structure that SDKs send
type Event struct {
	// Required fields
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
	Platform  string    `json:"platform"`

	// Standard fields
	Level       string `json:"level,omitempty"` // debug, info, warning, error, fatal
	Logger      string `json:"logger,omitempty"`
	Transaction string `json:"transaction,omitempty"`
	ServerName  string `json:"server_name,omitempty"`
	Release     string `json:"release,omitempty"`
	Dist        string `json:"dist,omitempty"`
	Environment string `json:"environment,omitempty"`
	Message     string `json:"message,omitempty"`

	// Structured data
	Tags        map[string]string      `json:"tags,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	Contexts    map[string]interface{} `json:"contexts,omitempty"`
	Modules     map[string]string      `json:"modules,omitempty"`
	Fingerprint []string               `json:"fingerprint,omitempty"`

	// User information
	User *User `json:"user,omitempty"`

	// Request information
	Request *Request `json:"request,omitempty"`

	// Exception chain
	Exception []Exception `json:"exception,omitempty"`

	// Breadcrumbs
	Breadcrumbs []*Breadcrumb `json:"breadcrumbs,omitempty"`

	// SDK information
	Sdk *SdkInfo `json:"sdk,omitempty"`

	// Transaction-specific (for performance monitoring)
	Type      string    `json:"type,omitempty"` // "event" or "transaction"
	StartTime time.Time `json:"start_timestamp,omitempty"`
	Spans     []*Span   `json:"spans,omitempty"`
}

// User describes the user associated with an event
type User struct {
	ID        string            `json:"id,omitempty"`
	Email     string            `json:"email,omitempty"`
	IPAddress string            `json:"ip_address,omitempty"`
	Username  string            `json:"username,omitempty"`
	Name      string            `json:"name,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}

// Request contains information about an HTTP request
type Request struct {
	URL         string            `json:"url,omitempty"`
	Method      string            `json:"method,omitempty"`
	Data        string            `json:"data,omitempty"`
	QueryString string            `json:"query_string,omitempty"`
	Cookies     string            `json:"cookies,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
}

// Exception represents an error in the exception chain
type Exception struct {
	Type       string      `json:"type,omitempty"`
	Value      string      `json:"value,omitempty"`
	Module     string      `json:"module,omitempty"`
	ThreadID   uint64      `json:"thread_id,omitempty"`
	Stacktrace *Stacktrace `json:"stacktrace,omitempty"`
	Mechanism  *Mechanism  `json:"mechanism,omitempty"`
}

// Stacktrace contains the frames of an exception
type Stacktrace struct {
	Frames []Frame `json:"frames,omitempty"`
}

// Frame represents a single stack frame
type Frame struct {
	Filename        string                 `json:"filename,omitempty"`
	Function        string                 `json:"function,omitempty"`
	Module          string                 `json:"module,omitempty"`
	Lineno          int                    `json:"lineno,omitempty"`
	Colno           int                    `json:"colno,omitempty"`
	AbsPath         string                 `json:"abs_path,omitempty"`
	ContextLine     string                 `json:"context_line,omitempty"`
	PreContext      []string               `json:"pre_context,omitempty"`
	PostContext     []string               `json:"post_context,omitempty"`
	InApp           bool                   `json:"in_app,omitempty"`
	Vars            map[string]interface{} `json:"vars,omitempty"`
	Package         string                 `json:"package,omitempty"`
	InstructionAddr string                 `json:"instruction_addr,omitempty"`
}

// Mechanism describes how an exception was generated
type Mechanism struct {
	Type        string                 `json:"type,omitempty"`
	Description string                 `json:"description,omitempty"`
	HelpLink    string                 `json:"help_link,omitempty"`
	Handled     *bool                  `json:"handled,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// Breadcrumb represents an event that occurred before the error
type Breadcrumb struct {
	Type      string                 `json:"type,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Level     string                 `json:"level,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

// Span represents a performance span (for transactions)
type Span struct {
	TraceID        string                 `json:"trace_id,omitempty"`
	SpanID         string                 `json:"span_id,omitempty"`
	ParentSpanID   string                 `json:"parent_span_id,omitempty"`
	Op             string                 `json:"op,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Status         string                 `json:"status,omitempty"`
	Tags           map[string]string      `json:"tags,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
	StartTimestamp time.Time              `json:"start_timestamp,omitempty"`
	Timestamp      time.Time              `json:"timestamp,omitempty"`
}

// ParseEvent parses a JSON payload into an Event
func ParseEvent(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// GetFingerprint returns the fingerprint for issue grouping
// If no fingerprint is provided, generates a default one
func (e *Event) GetFingerprint() []string {
	if len(e.Fingerprint) > 0 {
		return e.Fingerprint
	}

	// Default fingerprint based on exception type and message
	if len(e.Exception) > 0 {
		exc := e.Exception[0]
		return []string{exc.Type, exc.Value}
	}

	// Fallback to message
	if e.Message != "" {
		return []string{e.Message}
	}

	return []string{"default"}
}

// GetTitle returns a human-readable title for the event
func (e *Event) GetTitle() string {
	if len(e.Exception) > 0 {
		exc := e.Exception[0]
		if exc.Type != "" {
			return exc.Type
		}
		return exc.Value
	}
	if e.Message != "" {
		if len(e.Message) > 100 {
			return e.Message[:100] + "..."
		}
		return e.Message
	}
	return "Unknown Error"
}

// GetDescription returns a description for the event
func (e *Event) GetDescription() string {
	if len(e.Exception) > 0 {
		return e.Exception[0].Value
	}
	return e.Message
}

// ToJSON converts the event to JSON bytes
func (e *Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
