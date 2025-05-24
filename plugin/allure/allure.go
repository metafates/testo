package allure

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"testman/plugin"

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

func (a *Allure) Hooks() plugin.Hooks {
	return plugin.Hooks{
		BeforeAll: func() {
			fmt.Println("Allure.BeforeAll " + a.Name())
			a.start = time.Now()
		},
		BeforeEach: func() {
			fmt.Println("Allure.BeforeEach " + a.Name())
			a.start = time.Now()
		},
		AfterEach: func() {
			fmt.Println("Allure.AfterEach " + a.Name())
			a.stop = time.Now()
		},
		AfterAll: a.afterAll,
	}
}

func (a *Allure) Overrides() plugin.Overrides {
	return plugin.Overrides{
		Log: func(f plugin.FuncLog) plugin.FuncLog {
			return func(args ...any) {
				fmt.Println("inside log override " + a.Name())

				f(args...)
			}
		},
		Logf: func(f plugin.FuncLogf) plugin.FuncLogf {
			return func(format string, args ...any) {
				fmt.Println("inside logf override " + a.Name())

				f(format, args...)
			}
		},
	}
}

func (a *Allure) afterAll() {
	fmt.Println("Allure.AfterAll " + a.Name())
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

// Flaky indicates that this test or step is known to be unstable and can may not succeed every time.
func (a *Allure) Flaky() {
	a.statusDetails.Flaky = true
}

// Muted indicates that the result must not affect the statistics.
func (a *Allure) Muted() {
	a.statusDetails.Muted = true
}

// Known indicates that the test fails because of a known bug.
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
	return result{
		UUID:          uuid.NewString(),
		Name:          a.Name(),
		Links:         a.links,
		Labels:        a.labels,
		Status:        string(a.getStatus()),
		StatusDetails: a.statusDetails,
		Start:         a.start.UnixMilli(),
		Stop:          a.stop.UnixMilli(),
		Steps:         a.steps(),
	}
}

func (a *Allure) asStep() step {
	return step{
		Name:          a.BaseName(),
		Status:        a.getStatus(),
		StatusDetails: a.statusDetails,
		Start:         a.start.UnixMilli(),
		Stop:          a.stop.UnixMilli(),
		Steps:         a.steps(),
	}
}

func (a *Allure) steps() []step {
	steps := make([]step, 0, len(a.children))

	for _, c := range a.children {
		steps = append(steps, c.asStep())
	}

	return steps
}
