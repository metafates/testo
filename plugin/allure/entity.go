package allure

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

type result struct {
	UUID          string        `json:"uuid"`
	Name          string        `json:"name"`
	Links         []Link        `json:"links"`
	Labels        []Label       `json:"labels"`
	Status        string        `json:"status"`
	StatusDetails StatusDetails `json:"statusDetails"`
	Start         int64         `json:"start"`
	Stop          int64         `json:"stop"`
	Steps         []step        `json:"steps"`
}

type step struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
	Start  int64  `json:"start"`
	Stop   int64  `json:"stop"`
}
