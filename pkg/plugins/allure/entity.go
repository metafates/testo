package allure

import "github.com/google/uuid"

type stage int

const (
	stageTest stage = iota
	stageSetup
	stageTearDown
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
	Parameters    []Parameter   `json:"parameters,omitempty"`
}
