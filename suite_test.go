package testo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---- Tests for suiteCasesOf ----

type MySuite struct{}

func (MySuite) CasesNumbers() []int {
	return []int{1, 2, 3}
}

func (MySuite) CasesEmpty() []string {
	return []string{}
}

func TestSuiteCasesOf(t *testing.T) {
	cases := suiteCasesOf[MySuite](t)

	// Expect two case sets: "Numbers" and "Empty"
	assert.Contains(t, cases, "Numbers", "should find Numbers case set")
	assert.Contains(t, cases, "Empty", "should find Empty case set")

	// Test "Numbers" provides type int
	numCase := cases["Numbers"]
	assert.Equal(t, reflect.TypeOf(int(0)), numCase.Provides)

	// Run the case function
	numVals := numCase.Func(MySuite{})
	ints := make([]int, len(numVals))

	for i, rv := range numVals {
		ints[i] = int(rv.Int())
	}

	assert.Equal(t, []int{1, 2, 3}, ints, "Numbers case values should match")

	// Test "Empty" provides type string
	emptyCase := cases["Empty"]
	assert.Equal(t, reflect.TypeOf(""), emptyCase.Provides)

	// Run the case function
	emptyVals := emptyCase.Func(MySuite{})
	assert.Len(t, emptyVals, 0, "Empty case should yield zero values")
}

// ---- Tests for suiteHooksOf ----

var hookEvents []string

type HooksSuite struct{}

func (h *HooksSuite) BeforeAll(t *testing.T) {
	hookEvents = append(hookEvents, "BeforeAll")
}

func (h *HooksSuite) BeforeEach(t *testing.T) {
	hookEvents = append(hookEvents, "BeforeEach")
}

func (h *HooksSuite) AfterEach(t *testing.T) {
	hookEvents = append(hookEvents, "AfterEach")
}

func (h *HooksSuite) AfterAll(t *testing.T) {
	hookEvents = append(hookEvents, "AfterAll")
}

func TestSuiteHooksOf(t *testing.T) {
	hookEvents = nil

	hooks := suiteHooksOf[HooksSuite](t)

	var suite HooksSuite

	// Invoke hooks in order
	hooks.BeforeAll(suite, t)
	hooks.BeforeEach(suite, t)
	hooks.AfterEach(suite, t)
	hooks.AfterAll(suite, t)

	assert.Equal(t,
		[]string{"BeforeAll", "BeforeEach", "AfterEach", "AfterAll"},
		hookEvents,
		"hooks should fire in correct order",
	)
}

// ---- Tests for cloneSuite ----

type ClonerSuite struct {
	Value int
}

func (c ClonerSuite) Clone() ClonerSuite {
	return ClonerSuite{Value: c.Value + 100}
}

type PlainSuite struct {
	Value int
}

func TestCloneSuite_Cloner(t *testing.T) {
	orig := ClonerSuite{Value: 7}
	clone := cloneSuite(orig)

	// Clone should add 100
	assert.Equal(t, 107, clone.Value)
}

func TestCloneSuite_DeepClone(t *testing.T) {
	orig := PlainSuite{Value: 5}
	clone := cloneSuite(orig)

	// Modify original after cloning
	orig.Value = 999

	// clone.Value should remain the initial value
	assert.Equal(t, 5, clone.Value)
}
