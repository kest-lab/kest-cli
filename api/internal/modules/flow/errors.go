package flow

// FlowError is a transport-friendly error carrying the intended HTTP status.
type FlowError struct {
	Status  int
	Message string
}

func (e *FlowError) Error() string {
	return e.Message
}

func newFlowError(status int, message string) error {
	return &FlowError{
		Status:  status,
		Message: message,
	}
}
