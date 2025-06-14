package testo

import (
	"testing"

	"github.com/metafates/testo/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

type MockPluginWithNonPointerT struct{ T }

type InvalidT struct {
	*T

	MockPluginWithNonPointerT
}

func TestConstruct(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		options := []plugin.Option{
			{Value: "foo", Propagate: true},
			{Value: "bar", Propagate: false},
		}

		res := construct[MockT](t, nil, nil, options...)

		require.Equal(t, []plugin.Option{
			{Value: "foo", Propagate: true},
			{Value: "bar", Propagate: false},
		}, res.levelOptions)
		require.Equal(t, res.T, res.MockPluginWithT.T)

		require.NotEqual(t, res.Other, nil)
		require.Equal(t, res.T, res.Other.T)

		child := construct(t, &res, nil, plugin.Option{Value: "fizz"})

		require.Equal(t, res.T, child.T.parent)
		require.NotEqual(t, res, child)
		require.Equal(t, []plugin.Option{
			{Value: "fizz"},
			{Value: "foo", Propagate: true},
		}, child.MockPluginWithoutT.options)
	})

	t.Run("invalid", func(t *testing.T) {
		require.Panics(t, func() {
			construct[InvalidT](t, nil, nil)
		})
	})
}

func Test_casesPermutations(t *testing.T) {
	t.Run("multiple keys", func(t *testing.T) {
		input := map[string][]int{
			"A": {1, 2},
			"B": {3},
		}
		perms := casesPermutations(input)
		// Expect 2 permutations: {A:1,B:3} and {A:2,B:3}
		assert.Len(t, perms, 2)
		expected := []map[string]int{
			{"A": 1, "B": 3},
			{"A": 2, "B": 3},
		}
		// Convert each to comparable map for assertion
		for _, exp := range expected {
			assert.Contains(t, perms, exp)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		perms := casesPermutations(map[string][]string{})
		// Empty input yields exactly one empty map
		assert.Len(t, perms, 1)
		assert.Equal(t, map[string]string{} /* one empty map */, perms[0])
	})
}

type TestT struct {
	*T

	TestPlugin
}

type TestSuite struct {
	beforeAllTriggered bool
}

type TestPlugin struct{ *T }

func (t TestPlugin) Plugin() plugin.Spec {
	return plugin.Spec{
		Hooks: plugin.Hooks{
			BeforeAll: plugin.Hook{
				Func: func() {
					pluginBeforeAll = append(pluginBeforeAll, t.Name())
				},
			},
			BeforeEach: plugin.Hook{
				Func: func() {
					pluginBeforeEach = append(pluginBeforeEach, t.Name())
				},
			},
			AfterEach: plugin.Hook{
				Func: func() {
					pluginAfterEach = append(pluginAfterEach, t.Name())
				},
			},
			AfterAll: plugin.Hook{
				Func: func() {
					pluginAfterAll = append(pluginAfterAll, t.Name())
				},
			},
		},
	}
}

var (
	beforeAll  []string
	beforeEach []string
	afterEach  []string
	afterAll   []string

	pluginBeforeAll  []string
	pluginBeforeEach []string
	pluginAfterEach  []string
	pluginAfterAll   []string
)

func TestRunSuite(t *testing.T) {
	beforeAll = nil
	beforeEach = nil
	afterEach = nil
	afterAll = nil

	pluginBeforeAll = nil
	pluginBeforeEach = nil
	pluginAfterEach = nil
	pluginAfterAll = nil

	RunSuite[*TestSuite, *TestT](t)

	assert.Equal(t, []string{"TestRunSuite/TestSuite"}, beforeAll)
	assert.Equal(t, []string{"TestRunSuite/TestSuite"}, pluginBeforeAll)

	assert.Equal(t, []string{
		"TestRunSuite/TestSuite/Bar",
		"TestRunSuite/TestSuite/Foo",
	}, beforeEach)
	assert.Equal(t, []string{
		"TestRunSuite/TestSuite/Bar",
		"TestRunSuite/TestSuite/Foo",
		"TestRunSuite/TestSuite/Foo/subtest",
	}, pluginBeforeEach)

	assert.Equal(t, []string{
		"TestRunSuite/TestSuite/Bar",
		"TestRunSuite/TestSuite/Foo",
	}, afterEach)
	assert.Equal(t, []string{
		"TestRunSuite/TestSuite/Bar",
		"TestRunSuite/TestSuite/Foo/subtest",
		"TestRunSuite/TestSuite/Foo",
	}, pluginAfterEach)

	assert.Equal(t, []string{"TestRunSuite/TestSuite"}, afterAll)
	assert.Equal(t, []string{"TestRunSuite/TestSuite"}, pluginAfterAll)
}

func (s *TestSuite) BeforeAll(t *TestT) {
	s.beforeAllTriggered = true

	beforeAll = append(beforeAll, t.Name())
}

func (s TestSuite) BeforeEach(t *TestT) {
	assert.True(t, s.beforeAllTriggered)

	beforeEach = append(beforeEach, t.Name())
}

func (s TestSuite) AfterEach(t *TestT) {
	assert.True(t, s.beforeAllTriggered)

	afterEach = append(afterEach, t.Name())
}

func (s TestSuite) AfterAll(t *TestT) {
	assert.True(t, s.beforeAllTriggered)

	afterAll = append(afterAll, t.Name())
}

func (s TestSuite) TestFoo(t *TestT) {
	Run(t, "subtest", func(t *TestT) {})
}

func (s *TestSuite) TestBar(t *TestT) {}
