// Package allure provides Allure provider as a plugin for testo.
package allure

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/metafates/testo"
	"github.com/metafates/testo/internal/maputil"
	"github.com/metafates/testo/plugin"
)

// TODO: use tools.go pattern or go tool command when this plugin is moved into separate repo.

//go:generate ifacemaker -f $GOFILE -o interface.go -s Allure -i Interface -p $GOPACKAGE -e Plugin -e Init -y "Interface defines allure plugin interface. Useful for writing helpers which require allure methods but can't rely on concrete type."

var _ Interface = (*Allure)(nil)

// Allure defines allure plugin.
type Allure struct {
	*testo.T

	id UUID

	start, stop    time.Time
	rawLabels      []Label
	parameters     []Parameter
	links          []Link
	description    string
	rawStatus      Status
	statusDetails  StatusDetails
	categories     []Category
	rawAttachments []namedAttachment

	children []*Allure

	outputDir string
	stage     stage

	owner          string
	epic           string
	feature        string
	story          string
	severity       Severity
	titleOverwrite string
}

// Init plugin.
func (a *Allure) Init(parent *Allure, options ...plugin.Option) {
	a.id = uuid.NewString()
	a.outputDir = *outputDir

	for _, o := range options {
		if o, ok := o.Value.(option); ok {
			o(a)
		}
	}

	meta := testo.Inspect(a)

	info, ok := meta.Test.(plugin.RegularTestInfo)
	if ok {
		a.titleOverwrite = strings.TrimPrefix(info.RawBaseName, "Test")
	}

	if parent != nil {
		parent.children = append(parent.children, a)
	}
}

// Plugin implements [plugin.Plugin].
func (a *Allure) Plugin() plugin.Spec {
	return plugin.Spec{
		Hooks:     a.hooks(),
		Overrides: a.overrides(),
	}
}

// Title sets a human-readable title of the test.
//
// If not provided, function or subtest name is used instead.
func (a *Allure) Title(title string) {
	a.titleOverwrite = title
}

// Description sets an arbitrary text describing the test in
// more details than the title could fit.
//
// The description will be treated as a Markdown text,
// so you can you some basic formatting in it.
// HTML tags are not allowed in such a text and will
// be removed when building the report.
func (a *Allure) Description(desc string) {
	a.description = desc
}

// Links adds a list of links to webpages that may be useful for a reader investigating a test failure.
// You can provide as many links as needed.
//
// There are three types of links:
//   - a standard web link, e.g., a link to the description of the feature being tested;
//   - a link to an issue in the product's issue tracker;
//   - a link to the test description in a test management system (TMS).
func (a *Allure) Links(links ...Link) {
	a.links = append(a.links, links...)
}

// Status sets the status of the test.
func (a *Allure) Status(status Status) {
	a.rawStatus = status
}

// Labels adds given labels to the test result.
//
// A test result can have multiple labels with the same name.
// For example, this is often the case when a test result has multiple tags.
//
// Consider using helper methods such as [Allure.Tags] or [Allure.Severity]
// instead of using labels directly.
func (a *Allure) Labels(labels ...Label) {
	a.rawLabels = append(a.rawLabels, labels...)
}

// Tags adds short terms the test is related to.
// Usually it's a good idea to list relevant
// features that are being tested.
//
// Tags can then be used for [filtering].
//
// [filtering]: https://allurereport.org/docs/sorting-and-filtering/#filter-tests-by-tags
func (a *Allure) Tags(tags ...string) {
	for _, tag := range tags {
		a.rawLabels = append(a.rawLabels, Label{Name: "tag", Value: tag})
	}
}

// Parameters adds parameters to show for this report in the result.
func (a *Allure) Parameters(parameters ...Parameter) {
	a.parameters = append(a.parameters, parameters...)
}

// Owner sets the team member who is responsible for the test's stability.
// For example, this can be the test's author, the
// leading developer of the feature being tested, etc.
func (a *Allure) Owner(owner string) {
	a.owner = owner
}

// Severity sets a value indicating how important the test is.
// This may give the future reader an idea of how
// to prioritize the investigations of different test failures.
func (a *Allure) Severity(severity Severity) {
	a.severity = severity
}

// Epic linked to this test.
func (a *Allure) Epic(epic string) {
	a.epic = epic
}

// Feature linked to this test.
func (a *Allure) Feature(feature string) {
	a.feature = feature
}

// Story linked to this test.
func (a *Allure) Story(story string) {
	a.story = story
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

// Attach an attachment.
//
// See [NewAttachmentBytes] and [NewAttachmentPath] to create an attachment.
func (a *Allure) Attach(name string, attachment Attachment) {
	a.rawAttachments = append(a.rawAttachments, namedAttachment{
		Attachment: attachment,
		name:       name,
	})
}

func (a *Allure) panicked() bool {
	return testo.Inspect(a).Panic != nil
}

func (a *Allure) status() Status {
	if a.panicked() {
		return StatusBroken
	}

	if a.Skipped() {
		return StatusSkipped
	}

	if a.Failed() {
		return StatusFailed
	}

	return cmp.Or(a.rawStatus, StatusPassed)
}

func (a *Allure) asResult() result {
	return result{
		UUID:          a.id,
		FullName:      a.Name(),
		HistoryID:     a.Name(),
		Name:          a.title(),
		Links:         a.links,
		Parameters:    a.parameters,
		Labels:        a.labels(),
		Status:        a.status(),
		StatusDetails: a.statusDetails,
		Attachments:   a.attachments(),
		Start:         a.start.UnixMilli(),
		Stop:          a.stop.UnixMilli(),
		Steps:         a.steps(),
	}
}

func (a *Allure) attachments() []attachment {
	res := make([]attachment, 0, len(a.rawAttachments))

	for _, at := range a.rawAttachments {
		res = append(res, attachment{
			Name:   at.name,
			Source: filenameForAttachment(at),
			Type:   at.Type(),
		})
	}

	return res
}

func (a *Allure) allRawAttachments() []namedAttachment {
	attachments := slices.Clone(a.rawAttachments)

	for _, child := range a.children {
		attachments = append(attachments, child.allRawAttachments()...)
	}

	return attachments
}

func (a *Allure) title() string {
	return cmp.Or(
		a.titleOverwrite,
		strings.TrimPrefix(testBaseName(a.Name()), "Test"),
	)
}

func (a *Allure) asStep() step {
	return step{
		Name:          a.title(),
		Status:        a.status(),
		StatusDetails: a.statusDetails,
		Start:         a.start.UnixMilli(),
		Stop:          a.stop.UnixMilli(),
		Steps:         a.steps(),
		Parameters:    a.parameters,
		Attachments:   a.attachments(),
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
			UUID:     uuid.NewString(),
			Start:    start,
			Stop:     stop,
			Children: []UUID{child.id},
			Befores:  befores,
			Afters:   afters,
		})
	}

	return containers
}

func (a *Allure) beforeEach() {
	a.start = time.Now()

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
}

func (a *Allure) afterEach() {
	a.stop = time.Now()

	info := testo.Inspect(a)

	if info.Panic != nil {
		a.statusDetails.Message += fmt.Sprintf("panic: %v", info.Panic.Value)
		a.statusDetails.Trace = info.Panic.Trace
	}
}

func (a *Allure) hooks() plugin.Hooks {
	return plugin.Hooks{
		BeforeEach:    plugin.Hook{Func: a.beforeEach},
		BeforeEachSub: plugin.Hook{Func: a.beforeEach},
		AfterEachSub:  plugin.Hook{Func: a.afterEach},
		AfterEach:     plugin.Hook{Func: a.afterEach},
		AfterAll:      plugin.Hook{Func: a.afterAll},
	}
}

func (a *Allure) overrides() plugin.Overrides {
	return plugin.Overrides{
		Errorf: func(f plugin.FuncErrorf) plugin.FuncErrorf {
			return func(format string, args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message += trimLines(fmt.Sprintf(format, args...)) + "\n"

				f(format, args...)
			}
		},
		Error: func(f plugin.FuncError) plugin.FuncError {
			return func(args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message += trimLines(fmt.Sprint(args...)) + "\n"

				f(args...)
			}
		},
		Fatalf: func(f plugin.FuncFatalf) plugin.FuncFatalf {
			return func(format string, args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message = trimLines(fmt.Sprintf(format, args...)) + "\n"

				f(format, args...)
			}
		},
		Fatal: func(f plugin.FuncFatal) plugin.FuncFatal {
			return func(args ...any) {
				a.Helper()

				a.statusDetails.Trace = string(debug.Stack())
				a.statusDetails.Message = trimLines(fmt.Sprint(args...)) + "\n"

				f(args...)
			}
		},
	}
}

func (a *Allure) results() []result {
	results := make([]result, 0, len(a.children))

	type Name struct{ Full, Base string }

	parametrized := make(map[Name][]step)

	for _, test := range a.children {
		meta := testo.Inspect(test)

		if p, ok := meta.Test.(plugin.ParametrizedTestInfo); ok {
			name := Name{Full: removeCaseSuffix(test.Name()), Base: p.RawBaseName}

			parametrized[name] = append(parametrized[name], test.asStep())
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
			UUID: uuid.NewString(),
			Labels: []Label{
				{Name: "suite", Value: a.SuiteName()},
			},
			HistoryID: name.Full,
			FullName:  name.Full,
			Name:      name.Base,
			Status:    status,
			Start:     start,
			Stop:      stop,
			Steps:     steps,
		})
	}

	return results
}

func (a *Allure) afterAll() {
	err := os.Mkdir(a.outputDir, 0o750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		a.Fatal(err)
	}

	a.writeResults()
	a.writeContainers()
	a.writeCategories()
	a.writeAttachments()
	a.writeProperties()
}

func (a *Allure) writeResults() {
	for _, res := range a.results() {
		marshalled, err := json.Marshal(res)
		if err != nil {
			a.Fatalf("marshal: %v", err)
		}

		err = os.WriteFile(
			filepath.Join(a.outputDir, res.UUID+"-result.json"),
			marshalled,
			0o600,
		)
		if err != nil {
			a.Fatalf("write result: %v", err)
		}
	}
}

func (a *Allure) writeContainers() {
	for _, c := range a.containers() {
		marshalled, err := json.Marshal(c)
		if err != nil {
			a.Fatalf("marshal: %v", err)
		}

		err = os.WriteFile(
			filepath.Join(a.outputDir, c.UUID+"-container.json"),
			marshalled,
			0o600,
		)
		if err != nil {
			a.Fatalf("write container: %v", err)
		}
	}
}

func (a *Allure) writeAttachments() {
	for _, at := range a.allRawAttachments() {
		a.writeAttachment(at)
	}
}

func (a *Allure) writeAttachment(attachment Attachment) {
	src, err := attachment.Open()
	if err != nil {
		a.Fatalf("failed to open attachment: %v", err)
	}
	defer src.Close()

	path := a.attachmentPath(attachment)

	dst, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		a.Fatalf("failed to create file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		a.Fatalf("failed to copy files: %v", err)
	}
}

func (a *Allure) writeProperties() {
	p := newProperties()

	marshalled, err := p.MarshalProperties()
	if err != nil {
		a.Fatalf("marshal properties: %v", err)
	}

	err = os.WriteFile(filepath.Join(a.outputDir, "environment.properties"), marshalled, 0o600)
	if err != nil {
		a.Fatalf("write properties: %v", err)
	}
}

//nolint:gochecknoglobals // global variable is required here
var writeCategoriesMutex sync.Mutex

func (a *Allure) writeCategories() {
	// If multiple suites are run in parallel, there exists a small
	// chance that they will finish at the same time.
	// In that case categories file won't be written properly.
	writeCategoriesMutex.Lock()
	defer writeCategoriesMutex.Unlock()

	// This is tricky.
	// We could already have categories file written
	// by other suite, so we need to append to it.
	// But also we have to remain categories unique.
	path := filepath.Join(a.outputDir, "categories.json")

	readExisting := func() []Category {
		file, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}

			a.Fatalf("read categories: %v", err)
		}

		var out []Category

		// if json is malformed we should ignore it and overwrite.
		_ = json.Unmarshal(file, &out)

		return out
	}

	categories := readExisting()
	categories = append(categories, a.categories...)
	categories = uniqueCategories(categories)

	if len(categories) == 0 {
		return
	}

	marshalled, err := json.Marshal(categories)
	if err != nil {
		a.Fatalf("marshal category: %v", err)
	}

	err = os.WriteFile(path, marshalled, 0o600)
	if err != nil {
		a.Fatalf("write categories: %v", err)
	}
}

func (a *Allure) labels() []Label {
	labels := uniqueLabels(a.rawLabels)

	hostname, _ := os.Hostname()

	// TODO: restrict adding these labels from Labels method
	for _, l := range []Label{
		{Name: "suite", Value: a.SuiteName()},
		{Name: "owner", Value: a.owner},
		{Name: "epic", Value: a.epic},
		{Name: "feature", Value: a.feature},
		{Name: "story", Value: a.story},
		{Name: "severity", Value: string(a.severity)},
		{Name: "host", Value: hostname},
		{Name: "language", Value: "go"},
	} {
		if l.Value != "" {
			labels = append(labels, l)
		}
	}

	return labels
}

func (a *Allure) attachmentPath(attachment Attachment) string {
	return filepath.Join(a.outputDir, filenameForAttachment(attachment))
}

func newProperties() properties {
	return properties{
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		GoVersion: runtime.Version(),
	}
}

func trimLines(s string) string {
	s = strings.TrimSpace(s)

	lines := make([]string, 0, strings.Count(s, "\n"))

	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func uniqueCategories(categories []Category) []Category {
	byName := make(map[string]Category, len(categories))

	for _, l := range categories {
		byName[l.Name] = l
	}

	sortedKeys := maputil.Keys(byName)
	slices.Sort(sortedKeys)

	unique := make([]Category, 0, len(sortedKeys))

	for _, k := range sortedKeys {
		unique = append(unique, byName[k])
	}

	return unique
}

func uniqueLabels(labels []Label) []Label {
	byName := make(map[string]Label, len(labels))

	for _, l := range labels {
		byName[l.Name] = l
	}

	sortedKeys := maputil.Keys(byName)
	slices.Sort(sortedKeys)

	unique := make([]Label, 0, len(sortedKeys))

	for _, k := range sortedKeys {
		unique = append(unique, byName[k])
	}

	return unique
}

type namedAttachment struct {
	Attachment

	name string
}

func filenameForAttachment(attachment Attachment) string {
	byType, _ := mime.ExtensionsByType(attachment.Type().String())

	ext := ".txt"
	if len(byType) > 0 {
		ext = byType[0]
	}

	// {uuid}-attachment.{ext}
	return attachment.UUID() + "-attachment" + ext
}

func removeCaseSuffix(testName string) string {
	// RealName_case_%d
	//         ^
	idx := strings.LastIndex(testName, "_case_")

	if idx < 0 {
		return testName
	}

	return testName[:idx]
}

func testBaseName(testName string) string {
	segments := strings.Split(testName, "/")

	if len(segments) == 0 {
		return ""
	}

	return segments[len(segments)-1]
}
