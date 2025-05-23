package allure

import (
	"cmp"
	"encoding/json"
	"os"
	"time"

	"testman"

	"github.com/google/uuid"
)

type Allure struct {
	*testman.T

	start, stop time.Time

	labels        []Label
	links         []Link
	description   string
	status        Status
	statusDetails StatusDetails

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

func (a *Allure) BeforeAll() {
	a.start = time.Now()
}

func (a *Allure) AfterAll() {
	a.stop = time.Now()

	for _, test := range a.children {
		res := test.asResult()

		resJSON, _ := json.Marshal(res)

		os.Mkdir("allure-results", os.ModePerm)
		os.WriteFile("allure-results/"+res.UUID+"-result.json", resJSON, os.ModePerm)
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

func (a *Allure) Flaky() {
	a.statusDetails.Flaky = true
}

func (a *Allure) Muted() {
	a.statusDetails.Muted = true
}

func (a *Allure) Known() {
	a.statusDetails.Known = true
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

func (a *Allure) asResult() result {
	steps := make([]step, 0, len(a.children))

	for _, c := range a.children {
		steps = append(steps, c.asStep())
	}

	return result{
		UUID:   uuid.NewString(),
		Name:   a.Name(),
		Links:  nilAsEmpty(a.links),
		Labels: nilAsEmpty(a.labels),
		Status: string(a.getStatus()),
		Start:  a.start.UnixMilli(),
		Stop:   a.stop.UnixMilli(),
		Steps:  steps,
	}
}

func (a *Allure) asStep() step {
	return step{
		Name:   a.BaseName(),
		Status: a.getStatus(),
		Start:  a.start.UnixMilli(),
		Stop:   a.stop.UnixMilli(),
	}
}

func nilAsEmpty[S ~[]T, T any](s S) S {
	if s == nil {
		return make(S, 0)
	}

	return s
}
