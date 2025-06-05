package allure

import "github.com/google/uuid"

type stage int

const (
	stageTest stage = iota
	stageSetup
	stageTearDown
)

type Severity string

const (
	SeverityTrivial  Severity = "trivial"
	SeverityMinor    Severity = "minor"
	SeverityNormal   Severity = "normal"
	SeverityCritical Severity = "critical"
	SeverityBlocker  Severity = "blocker"
)

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

type Parameter struct {
	Name     string        `json:"name"`
	Value    string        `json:"value"`
	Excluded bool          `json:"excluded"`
	Mode     ParameterMode `json:"mode"`
}

type ParameterMode string

const (
	ParameterModeDefault ParameterMode = "default"
	ParameterModeMasked  ParameterMode = "masked"
	ParameterModeHidden  ParameterMode = "hidden"
)

type Status string

const (
	StatusFailed  Status = "failed"
	StatusBroken  Status = "broken"
	StatusPassed  Status = "passed"
	StatusSkipped Status = "skipped"
	StatusUnknown Status = "unknown"
)

type Attachment struct {
	// Name is the human-readable name of the attachment.
	Name string `json:"name"`

	// Source is the name of the file with the attachment's content.
	Source string `json:"source"`

	// Type is the media type of the content.
	Type string `json:"type"`
}

type StatusDetails struct {
	Known   bool   `json:"known"`
	Muted   bool   `json:"muted"`
	Flaky   bool   `json:"flaky"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

type container struct {
	UUID     uuid.UUID  `json:"uuid"`
	Start    int64      `json:"start"`
	Stop     int64      `json:"stop"`
	Children uuid.UUIDs `json:"children"`
	Befores  []step     `json:"befores"`
	Afters   []step     `json:"afters"`
}

type result struct {
	UUID          uuid.UUID     `json:"uuid"`
	HistoryID     string        `json:"historyId"`
	FullName      string        `json:"fullName"`
	Name          string        `json:"name"`
	Links         []Link        `json:"links,omitempty"`
	Labels        []Label       `json:"labels,omitempty"`
	Parameters    []Parameter   `json:"parameters,omitempty"`
	Attachments   []Attachment  `json:"attachments,omitempty"`
	Status        Status        `json:"status"`
	StatusDetails StatusDetails `json:"statusDetails"`
	Start         int64         `json:"start"`
	Stop          int64         `json:"stop"`
	Steps         []step        `json:"steps,omitempty"`
}

type step struct {
	Name          string        `json:"name"`
	Status        Status        `json:"status"`
	StatusDetails StatusDetails `json:"statusDetails"`
	Start         int64         `json:"start"`
	Stop          int64         `json:"stop"`
	Steps         []step        `json:"steps,omitempty"`
	Attachments   []Attachment  `json:"attachments,omitempty"`
	Parameters    []Parameter   `json:"parameters,omitempty"`
}
