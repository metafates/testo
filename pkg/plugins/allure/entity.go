package allure

import (
	"bytes"

	"github.com/google/uuid"
)

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

type Category struct {
	// Name of the category.
	Name string `json:"name"`

	// MessageRegex is the regular expression
	// that the test result's message should match.
	MessageRegex string `json:"messageRegex"`

	// TraceRegex is the regular expression that
	// the test result's trace should match.
	TraceRegex string `json:"traceRegex"`

	// MatchedStatuses are the statuses that
	// the test result should be one of.
	MatchedStatuses []Status `json:"matchedStatuses"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LinkType string

const (
	LinkTypeIssue LinkType = "issue"
	LinkTypeTMS   LinkType = "tms"
)

type Link struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Type LinkType `json:"type"`
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

// TODO: use something like
// https://github.com/matishsiao/goInfo to extend properties.
type properties struct {
	OSPlatform string
	OSArch     string
	GoVersion  string
}

func (p properties) MarshalProperties() ([]byte, error) {
	var buf bytes.Buffer

	for _, line := range []struct{ Key, Value string }{
		{Key: "os_platform", Value: p.OSPlatform},
		{Key: "os_arch", Value: p.OSArch},
		{Key: "go_version", Value: p.GoVersion},
	} {
		_, err := buf.WriteString(line.Key + " = " + line.Value + "\n")
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
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
