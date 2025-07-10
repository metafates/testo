package allure

import (
	"bytes"
	"fmt"
)

// UUID is unique identifier.
type UUID = string

type stage int

const (
	stageTest stage = iota
	stageSetup
	stageTearDown
)

// Severity is a value indicating how important the test is.
// This may give the future reader an idea of how to prioritize
// the investigations of different test failures.
type Severity string

// Possible severity values.
const (
	SeverityTrivial  Severity = "trivial"
	SeverityMinor    Severity = "minor"
	SeverityNormal   Severity = "normal"
	SeverityCritical Severity = "critical"
	SeverityBlocker  Severity = "blocker"
)

// Category defines tests category.
//
// Allure checks each test against all the categories in the file,
// from top to bottom. The test gets assigned the first matching category.
// When no matches are found, Allure uses one of the default categories
// if the test is unsuccessful or no category otherwise.
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

// Label holds additional information about the test.
type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// LinkType is the type of link.
type LinkType string

// Possible link types.
const (
	LinkTypeIssue LinkType = "issue"
	LinkTypeTMS   LinkType = "tms"
)

// Link to webpage that may be useful for a reader
// investigating a test failure.
type Link struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Type LinkType `json:"type"`
}

// Parameter to show in the report.
//
// Allure plugin automatically sets parameters
// for parametrized tests.
type Parameter struct {
	Name     string        `json:"name"`
	Value    string        `json:"value"`
	Excluded bool          `json:"excluded"`
	Mode     ParameterMode `json:"mode"`
}

// NewParameter constructs a new [Parameter].
func NewParameter(name string, value any) Parameter {
	return Parameter{
		Name:  name,
		Value: fmt.Sprint(value),
		Mode:  ParameterModeDefault,
	}
}

// Masked returns a new parameter with mode set to masked.
func (p Parameter) Masked() Parameter {
	p.Mode = ParameterModeMasked

	return p
}

// Hidden returns a new parameter with mode set to hidden.
func (p Parameter) Hidden() Parameter {
	p.Mode = ParameterModeHidden

	return p
}

// ParameterMode controls how the parameter will be shown in the report.
type ParameterMode string

const (
	// ParameterModeDefault - the parameter and its value
	// will be shown in a table along with other parameters.
	ParameterModeDefault ParameterMode = "default"

	// ParameterModeMasked - the parameter will be shown
	// in the table, but its value will be hidden.
	ParameterModeMasked ParameterMode = "masked"

	// ParameterModeHidden - the parameter and its value
	// will not be shown in the test report.
	ParameterModeHidden ParameterMode = "hidden"
)

// Status is the test status.
type Status string

// Possible statuses.
const (
	StatusFailed  Status = "failed"
	StatusBroken  Status = "broken"
	StatusPassed  Status = "passed"
	StatusSkipped Status = "skipped"
	StatusUnknown Status = "unknown"
)

// StatusDetails holds additional information for status.
type StatusDetails struct {
	// Known indicates that the test
	// fails because of a known bug.
	Known bool `json:"known"`

	// Muted indicates that the result
	// must not affect the statistics.
	Muted bool `json:"muted"`

	// Flaky indicates that this test or step is known
	// to be unstable and can may not succeed every time.
	Flaky bool `json:"flaky"`

	// Message is the short text message to display in the
	// test details, such as a name of the exception that led to a failure.
	Message string `json:"message"`

	// Trace is the full stack trace to display in the test details.
	Trace string `json:"trace"`
}

type attachment struct {
	Name   string    `json:"name"`
	Source string    `json:"source"`
	Type   MediaType `json:"type"`
}

// TODO: use something like
// https://github.com/matishsiao/goInfo to extend properties.
type properties struct {
	GoOS      string
	GoArch    string
	GoVersion string
}

// MarshalProperties marshals this structure into [.properties] format.
//
// [.properties]: https://en.wikipedia.org/wiki/.properties
func (p properties) MarshalProperties() ([]byte, error) {
	var buf bytes.Buffer

	for _, line := range []struct{ Key, Value string }{
		{Key: "go_os", Value: p.GoOS},
		{Key: "go_arch", Value: p.GoArch},
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
	UUID     UUID   `json:"uuid"`
	Start    int64  `json:"start"`
	Stop     int64  `json:"stop"`
	Children []UUID `json:"children"`
	Befores  []step `json:"befores"`
	Afters   []step `json:"afters"`
}

type result struct {
	UUID          UUID          `json:"uuid"`
	HistoryID     string        `json:"historyId"`
	FullName      string        `json:"fullName"`
	Name          string        `json:"name"`
	Links         []Link        `json:"links,omitempty"`
	Labels        []Label       `json:"labels,omitempty"`
	Parameters    []Parameter   `json:"parameters,omitempty"`
	Attachments   []attachment  `json:"attachments,omitempty"`
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
	Attachments   []attachment  `json:"attachments,omitempty"`
	Parameters    []Parameter   `json:"parameters,omitempty"`
}
