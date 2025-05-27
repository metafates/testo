package allure

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"testman"
	"testman/plugin"
)

type Allure struct {
	*testman.T

	start, stop time.Time

	labels        []Label
	parameters    []Parameter
	links         []Link
	description   string
	status        Status
	statusDetails StatusDetails

	children []*Allure

	// an example of the field set through options
	outputPath string
}

func (a *Allure) Init(parent *Allure, options ...plugin.Option) {
	a.outputPath = "allure-results"

	for _, o := range options {
		if o, ok := o.Value.(Option); ok {
			o(a)
		}
	}

	if parent != nil {
		parent.children = append(parent.children, a)
	}
}

func (a *Allure) Plugin() plugin.Plugin {
	return plugin.Plugin{
		Hooks:     a.hooks(),
		Overrides: a.overrides(),
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
	if a.Panicked() {
		return StatusBroken
	}

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
		HistoryID:     a.Name(),
		Name:          a.Name(),
		Links:         a.links,
		Parameters:    a.parameters,
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
		Parameters:    a.parameters,
	}
}

func (a *Allure) steps() []step {
	steps := make([]step, 0, len(a.children))

	for _, c := range a.children {
		steps = append(steps, c.asStep())
	}

	return steps
}

func (a *Allure) hooks() plugin.Hooks {
	return plugin.Hooks{
		BeforeAll: plugin.Hook{
			Func: func() {
				a.Log("Allure.BeforeAll " + a.Name())

				a.start = time.Now()
			},
		},
		BeforeEach: plugin.Hook{
			Func: func() {
				a.Log("Allure.BeforeEach " + a.Name())

				a.start = time.Now()
				a.labels = append(
					a.labels,
					Label{Name: "suite", Value: a.SuiteName()},
				)

				for name, value := range a.CaseParams() {
					a.parameters = append(a.parameters, Parameter{
						Name:  name,
						Value: fmt.Sprint(value),
						Mode:  ParameterModeDefault,
					})
				}
			},
		},
		AfterEach: plugin.Hook{
			Func: func() {
				a.Log("Allure.AfterEach " + a.Name())

				a.stop = time.Now()

				if info, ok := a.PanicInfo(); ok {
					a.statusDetails.Message += fmt.Sprintf("panic: %v", info)
					a.statusDetails.Trace = string(debug.Stack())
				}
			},
		},
		AfterAll: plugin.Hook{Func: a.afterAll},
	}
}

func (a *Allure) overrides() plugin.Overrides {
	return plugin.Overrides{
		Log: func(f plugin.FuncLog) plugin.FuncLog {
			return func(args ...any) {
				a.Helper()

				fmt.Println("inside log override " + a.Name())

				f(args...)
			}
		},
		Logf: func(f plugin.FuncLogf) plugin.FuncLogf {
			return func(format string, args ...any) {
				a.Helper()

				fmt.Println("inside logf override " + a.Name())

				f(format, args...)
			}
		},
		Errorf: func(f plugin.FuncErrorf) plugin.FuncErrorf {
			return func(format string, args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message += fmt.Sprintf(format, args...) + "\n"
				f(format, args...)
			}
		},
		Error: func(f plugin.FuncError) plugin.FuncError {
			return func(args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message += fmt.Sprint(args...) + "\n"
				f(args...)
			}
		},
		Fatalf: func(f plugin.FuncFatalf) plugin.FuncFatalf {
			return func(format string, args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message += fmt.Sprintf(format, args...) + "\n"
				f(format, args...)
			}
		},
		Fatal: func(f plugin.FuncFatal) plugin.FuncFatal {
			return func(args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message += fmt.Sprint(args...) + "\n"
				f(args...)
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

		os.Mkdir(a.outputPath, os.ModePerm)
		os.WriteFile(filepath.Join(a.outputPath, res.UUID+"-result.json"), resJSON, os.ModePerm)
	}
}
