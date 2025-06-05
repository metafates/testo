// Package allure provides Allure provider as a plugin for testo.
package allure

import (
	"cmp"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/metafates/testo"
	"github.com/metafates/testo/plugin"
)

var outputDir = flag.String(
	"allure.output",
	"allure-results",
	"path to output dir for allure results",
)

type Allure struct {
	*testo.T

	start, stop time.Time

	labels        []Label
	parameters    []Parameter
	links         []Link
	description   string
	status        Status
	statusDetails StatusDetails

	children []*Allure

	outputPath string
	stage      stage

	id uuid.UUID
}

func (a *Allure) Init(parent *Allure, options ...plugin.Option) {
	a.id = uuid.New()
	a.outputPath = *outputDir

	for _, o := range options {
		if o, ok := o.Value.(Option); ok {
			o(a)
		}
	}

	if parent != nil {
		parent.children = append(parent.children, a)
	}
}

func (a *Allure) Plugin() plugin.Spec {
	return plugin.Spec{
		Hooks:     a.hooks(),
		Overrides: a.overrides(),
	}
}

func (a *Allure) Title(title string) {
	// TODO
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

func (a *Allure) Tags(tags ...string) {
	a.Labels(newLabels("tag", tags...)...)
}

func (a *Allure) Owners(owners ...string) {
	a.Labels(newLabels("owner", owners...)...)
}

func (a *Allure) Severity(severity Severity) {
	a.Labels(Label{Name: "severity", Value: string(severity)})
}

// Flaky indicates that this test or step is known
// to be unstable and can may not succeed every time.
func (a *Allure) Flaky() {
	a.statusDetails.Flaky = true
}

// Muted indicates that the result
// must not affect the statistics.
func (a *Allure) Muted() {
	a.statusDetails.Muted = true
}

// Known indicates that the test
// fails because of a known bug.
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
		UUID:          a.id,
		FullName:      a.Name(),
		HistoryID:     a.Name(),
		Name:          a.BaseName(),
		Links:         a.links,
		Parameters:    a.parameters,
		Labels:        a.labels,
		Status:        a.getStatus(),
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
		if c.stage == stageTest {
			steps = append(steps, c.asStep())
		}
	}

	return steps
}

func (a *Allure) containers() []container {
	containers := make([]container, 0, len(a.children))

	for _, child := range a.children {
		var befores, afters []step

		var start, stop int64

		for _, sub := range child.children {
			switch sub.stage {
			case stageSetup:
				s := sub.asStep()

				befores = append(befores, s)
				start = min(start, s.Start)

			case stageTearDown:
				s := sub.asStep()

				afters = append(afters, s)
				stop = max(stop, s.Stop)
			}
		}

		containers = append(containers, container{
			UUID:     uuid.New(),
			Start:    start,
			Stop:     stop,
			Children: uuid.UUIDs{child.id},
			Befores:  befores,
			Afters:   afters,
		})
	}

	return containers
}

func (a *Allure) hooks() plugin.Hooks {
	return plugin.Hooks{
		BeforeEach: plugin.Hook{
			Func: func() {
				a.start = time.Now()
				a.labels = append(
					a.labels,
					Label{Name: "suite", Value: a.SuiteName()},
				)

				meta := testo.Inspect(a)

				if p, ok := meta.Test.(plugin.ParametrizedTestInfo); ok {
					for name, value := range p.Params {
						a.parameters = append(a.parameters, Parameter{
							Name:  name,
							Value: fmt.Sprint(value),
							Mode:  ParameterModeDefault,
						})
					}
				}
			},
		},
		AfterEach: plugin.Hook{
			Func: func() {
				a.stop = time.Now()

				if info, ok := a.PanicInfo(); ok {
					a.statusDetails.Message += fmt.Sprintf("panic: %v", info.Value)
					a.statusDetails.Trace = info.Trace
				}
			},
		},
		AfterAll: plugin.Hook{Func: a.afterAll},
	}
}

func (a *Allure) overrides() plugin.Overrides {
	return plugin.Overrides{
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
				a.statusDetails.Message = fmt.Sprintf(format, args...) + "\n"
				f(format, args...)
			}
		},
		Fatal: func(f plugin.FuncFatal) plugin.FuncFatal {
			return func(args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message = fmt.Sprint(args...) + "\n"
				f(args...)
			}
		},
	}
}

func (a *Allure) results() []result {
	results := make([]result, 0, len(a.children))

	parametrized := make(map[string][]step)

	for _, test := range a.children {
		meta := testo.Inspect(test)

		if p, ok := meta.Test.(plugin.ParametrizedTestInfo); ok {
			parametrized[p.BaseName] = append(parametrized[p.BaseName], test.asStep())
		} else {
			results = append(results, test.asResult())
		}
	}

	for name, steps := range parametrized {
		status := StatusPassed

		start := int64(math.MaxInt64)
		stop := int64(math.MinInt64)

		for _, s := range steps {
			switch s.Status {
			case StatusFailed, StatusBroken:
				status = s.Status
			}

			start = min(start, s.Start)
			stop = max(stop, s.Stop)
		}

		results = append(results, result{
			UUID: uuid.New(),
			Labels: []Label{
				{Name: "suite", Value: a.SuiteName()},
			},
			HistoryID: name,
			FullName:  name,
			Name:      name,
			Status:    status,
			Start:     start,
			Stop:      stop,
			Steps:     steps,
		})
	}

	return results
}

func (a *Allure) afterAll() {
	if len(a.children) > 0 {
		err := os.Mkdir(a.outputPath, 0o750)
		if err != nil && !errors.Is(err, os.ErrExist) {
			a.Fatal(err)
		}
	}

	for _, res := range a.results() {
		marshalled, err := json.Marshal(res)
		if err != nil {
			a.Fatalf("marshal: %v", err)
		}

		err = os.WriteFile(
			filepath.Join(a.outputPath, res.UUID.String()+"-result.json"),
			marshalled,
			0o600,
		)
		if err != nil {
			a.Fatalf("write file: %v", err)
		}
	}

	for _, c := range a.containers() {
		marshalled, err := json.Marshal(c)
		if err != nil {
			a.Fatalf("marshal: %v", err)
		}

		err = os.WriteFile(
			filepath.Join(a.outputPath, c.UUID.String()+"-container.json"),
			marshalled,
			0o600,
		)
		if err != nil {
			a.Fatalf("write file: %v", err)
		}
	}
}

func newLabels(name string, values ...string) []Label {
	labels := make([]Label, 0, len(values))

	for _, v := range values {
		labels = append(labels, Label{Name: name, Value: v})
	}

	return labels
}
