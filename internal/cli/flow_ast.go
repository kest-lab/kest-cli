package cli

type FlowMeta struct {
	ID      string
	Name    string
	Version string
	Env     string
	Tags    []string
}

type FlowStep struct {
	ID          string
	Name        string
	Type        string
	Retry       int
	RetryWait   int
	MaxDuration int
	OnFail      string
	LineNum     int
	Raw         string
	Request     RequestOptions
	Exec        ExecOptions
}

type ExecOptions struct {
	Command  string
	Captures []string
}

type FlowEdge struct {
	From    string
	To      string
	On      string
	LineNum int
}

type FlowDoc struct {
	Meta  FlowMeta
	Steps []FlowStep
	Edges []FlowEdge
}
