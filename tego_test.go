package tego

import (
	"reflect"
	"testing"

	"github.com/metafates/tego/plugin"
)

type MockT struct {
	*T

	MockPluginWithT
	MockPluginWithoutT
	Other *MockPluginWithT
}

type MockPluginWithT struct{ *T }

type MockPluginWithoutT struct {
	parent  *MockPluginWithoutT
	options []plugin.Option
}

func (m *MockPluginWithoutT) Init(parent *MockPluginWithoutT, options ...plugin.Option) {
	m.parent = parent
	m.options = options
}

func TestConstruct(t *testing.T) {
	options := []plugin.Option{
		{Value: "foo"},
		{Value: "bar"},
	}

	res := construct[MockT](t, nil, options...)

	equal(t, []plugin.Option{
		{Value: "foo"},
		{Value: "bar"},
	}, res.levelOptions)
	equal(t, res.T, res.MockPluginWithT.T)

	notEqual(t, res.Other, nil)
	equal(t, res.T, res.Other.T)

	child := construct(t, &res, plugin.Option{Value: "fizz"})

	equal(t, res.T, child.T.parent)
	notEqual(t, res, child)
	equal(t, []plugin.Option{
		{Value: "fizz"},
		{Value: "foo"},
		{Value: "bar"},
	}, child.MockPluginWithoutT.options)
}

func notEqual[T any](t *testing.T, a, b T) {
	t.Helper()

	if reflect.DeepEqual(a, b) {
		t.Fatalf("must not be equal: a %v, b %v", a, b)
	}
}

func equal[T any](t *testing.T, want, got T) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("must be equal: want %v, got %v", want, got)
	}
}
