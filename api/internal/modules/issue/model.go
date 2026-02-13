package issue

import "time"

// IssueStatus represents the status of an issue
type IssueStatus string

const (
	IssueStatusUnresolved IssueStatus = "unresolved"
	IssueStatusResolved   IssueStatus = "resolved"
	IssueStatusIgnored    IssueStatus = "ignored"
)

// Issue represents an aggregated problem (group of similar events)
type Issue struct {
	ProjectID   uint64
	Fingerprint string

	FirstSeen  time.Time
	LastSeen   time.Time
	EventCount uint64

	LastMessage string
	LastLevel   string
	LastEventID string

	Status IssueStatus
}

// IsResolved returns whether the issue is resolved
func (i *Issue) IsResolved() bool {
	return i.Status == IssueStatusResolved
}

// IsIgnored returns whether the issue is ignored
func (i *Issue) IsIgnored() bool {
	return i.Status == IssueStatusIgnored
}

// IsActive returns whether the issue is active (not resolved or ignored)
func (i *Issue) IsActive() bool {
	return i.Status == IssueStatusUnresolved
}
