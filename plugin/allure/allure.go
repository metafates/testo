package allure

import (
	"cmp"
	"fmt"
	"time"

	"testman"
)

type Allure struct {
	*testman.T

	start, stop time.Time

	labels      []Label
	links       []Link
	description string
	status      Status

	children []*Allure
}

func (a *Allure) New(t *testman.T) *Allure {
	child := &Allure{T: t}

	if a != nil {
		a.children = append(a.children, child)
	}

	return child
}

func (a *Allure) BeforeEach() {
	a.start = time.Now()
}

func (a *Allure) AfterEach() {
	a.stop = time.Now()
}

func (a *Allure) AfterAll() {
	fmt.Println("took", a.stop.Sub(a.start))

	for _, child := range a.children {
		fmt.Println("child", child.Name(), child.getStatus())
	}
}

func (a *Allure) Description(desc string) {
	a.description = desc
}

func (a *Allure) Links(links ...Link) {
	a.links = append(a.links, links...)
}

func (a *Allure) Status(status Status) {
	a.status = status
}

func (a *Allure) Labels(labels ...Label) {
	a.labels = append(a.labels, labels...)
}

func (a *Allure) getStatus() Status {
	if a.Skipped() {
		return StatusSkipped
	}

	if a.Failed() {
		return StatusFailed
	}

	return cmp.Or(a.status, StatusPassed)
}

type Label struct {
	Name  string
	Value string
}

type Link struct {
	Name string
	URL  string
	Type string
}

type Parameter struct {
	Name  string
	Value string
	Mode  ParameterMode
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
