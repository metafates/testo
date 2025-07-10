// Code generated from testify a53be35c3b0cfcd5189cffcfd75df60ea581104c; DO NOT EDIT.

package allure

import (
	"cmp"
	testo "github.com/metafates/testo"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// Condition uses a Comparison to assert a complex condition.
func (x Requirements[T]) Condition(comp assert.Comparison, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "condition")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("comp", comp).withMode(x.mode))
		require.Condition(t, comp, msgAndArgs...)
	})
}

// Condition uses a Comparison to assert a complex condition.
func (x Assertions[T]) Condition(comp assert.Comparison, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "condition")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("comp", comp).withMode(x.mode))
		assert.Condition(t, comp, msgAndArgs...)
	})
}

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
// 	assert.Contains(t, "Hello World", "World")
// 	assert.Contains(t, ["Hello", "World"], "World")
// 	assert.Contains(t, {"Hello": "World"}, "Hello")
func (x Requirements[T]) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "contains")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("s", s).withMode(x.mode), NewParameter("contains", contains).withMode(x.mode))
		require.Contains(t, s, contains, msgAndArgs...)
	})
}

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
// 	assert.Contains(t, "Hello World", "World")
// 	assert.Contains(t, ["Hello", "World"], "World")
// 	assert.Contains(t, {"Hello": "World"}, "Hello")
func (x Assertions[T]) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "contains")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("s", s).withMode(x.mode), NewParameter("contains", contains).withMode(x.mode))
		assert.Contains(t, s, contains, msgAndArgs...)
	})
}

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func (x Requirements[T]) DirExists(path string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "dir exists")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		require.DirExists(t, path, msgAndArgs...)
	})
}

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func (x Assertions[T]) DirExists(path string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "dir exists")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		assert.DirExists(t, path, msgAndArgs...)
	})
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func (x Requirements[T]) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "elements match")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list A", listA).withMode(x.mode), NewParameter("list B", listB).withMode(x.mode))
		require.ElementsMatch(t, listA, listB, msgAndArgs...)
	})
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// assert.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func (x Assertions[T]) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "elements match")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list A", listA).withMode(x.mode), NewParameter("list B", listB).withMode(x.mode))
		assert.ElementsMatch(t, listA, listB, msgAndArgs...)
	})
}

// Empty asserts that the given value is "empty".
//
// [Zero values] are "empty".
//
// Arrays are "empty" if every element is the zero value of the type (stricter than "empty").
//
// Slices, maps and channels with zero length are "empty".
//
// Pointer values are "empty" if the pointer is nil or if the pointed value is "empty".
//
// 	assert.Empty(t, obj)
//
// [Zero values]: https://go.dev/ref/spec#The_zero_value
func (x Requirements[T]) Empty(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "empty")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.Empty(t, object, msgAndArgs...)
	})
}

// Empty asserts that the given value is "empty".
//
// [Zero values] are "empty".
//
// Arrays are "empty" if every element is the zero value of the type (stricter than "empty").
//
// Slices, maps and channels with zero length are "empty".
//
// Pointer values are "empty" if the pointer is nil or if the pointed value is "empty".
//
// 	assert.Empty(t, obj)
//
// [Zero values]: https://go.dev/ref/spec#The_zero_value
func (x Assertions[T]) Empty(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "empty")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.Empty(t, object, msgAndArgs...)
	})
}

// Equal asserts that two objects are equal.
//
// 	assert.Equal(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func (x Requirements[T]) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.Equal(t, expected, actual, msgAndArgs...)
	})
}

// Equal asserts that two objects are equal.
//
// 	assert.Equal(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func (x Assertions[T]) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.Equal(t, expected, actual, msgAndArgs...)
	})
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
// 	actualObj, err := SomeFunction()
// 	assert.EqualError(t, err,  expectedErrorString)
func (x Requirements[T]) EqualError(theError error, errString string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal error")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("the error", theError).withMode(x.mode), NewParameter("err string", errString).withMode(x.mode))
		require.EqualError(t, theError, errString, msgAndArgs...)
	})
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
// 	actualObj, err := SomeFunction()
// 	assert.EqualError(t, err,  expectedErrorString)
func (x Assertions[T]) EqualError(theError error, errString string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal error")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("the error", theError).withMode(x.mode), NewParameter("err string", errString).withMode(x.mode))
		assert.EqualError(t, theError, errString, msgAndArgs...)
	})
}

// EqualExportedValues asserts that the types of two objects are equal and their public
// fields are also equal. This is useful for comparing structs that have private fields
// that could potentially differ.
//
// 	 type S struct {
// 		Exported     	int
// 		notExported   	int
// 	 }
// 	 assert.EqualExportedValues(t, S{1, 2}, S{1, 3}) => true
// 	 assert.EqualExportedValues(t, S{1, 2}, S{2, 3}) => false
func (x Requirements[T]) EqualExportedValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal exported values")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.EqualExportedValues(t, expected, actual, msgAndArgs...)
	})
}

// EqualExportedValues asserts that the types of two objects are equal and their public
// fields are also equal. This is useful for comparing structs that have private fields
// that could potentially differ.
//
// 	 type S struct {
// 		Exported     	int
// 		notExported   	int
// 	 }
// 	 assert.EqualExportedValues(t, S{1, 2}, S{1, 3}) => true
// 	 assert.EqualExportedValues(t, S{1, 2}, S{2, 3}) => false
func (x Assertions[T]) EqualExportedValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal exported values")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.EqualExportedValues(t, expected, actual, msgAndArgs...)
	})
}

// EqualValues asserts that two objects are equal or convertible to the larger
// type and equal.
//
// 	assert.EqualValues(t, uint32(123), int32(123))
func (x Requirements[T]) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal values")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.EqualValues(t, expected, actual, msgAndArgs...)
	})
}

// EqualValues asserts that two objects are equal or convertible to the larger
// type and equal.
//
// 	assert.EqualValues(t, uint32(123), int32(123))
func (x Assertions[T]) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "equal values")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.EqualValues(t, expected, actual, msgAndArgs...)
	})
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
// 	actualObj, err := SomeFunction()
// 	assert.Error(t, err)
func (x Requirements[T]) Error(err error, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode))
		require.Error(t, err, msgAndArgs...)
	})
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
// 	actualObj, err := SomeFunction()
// 	assert.Error(t, err)
func (x Assertions[T]) Error(err error, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode))
		assert.Error(t, err, msgAndArgs...)
	})
}

// ErrorAs asserts that at least one of the errors in err's chain matches target, and if so, sets target to that error value.
// This is a wrapper for errors.As.
func (x Requirements[T]) ErrorAs(err error, target interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error as")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		require.ErrorAs(t, err, target, msgAndArgs...)
	})
}

// ErrorAs asserts that at least one of the errors in err's chain matches target, and if so, sets target to that error value.
// This is a wrapper for errors.As.
func (x Assertions[T]) ErrorAs(err error, target interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error as")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		assert.ErrorAs(t, err, target, msgAndArgs...)
	})
}

// ErrorContains asserts that a function returned an error (i.e. not `nil`)
// and that the error contains the specified substring.
//
// 	actualObj, err := SomeFunction()
// 	assert.ErrorContains(t, err,  expectedErrorSubString)
func (x Requirements[T]) ErrorContains(theError error, contains string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error contains")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("the error", theError).withMode(x.mode), NewParameter("contains", contains).withMode(x.mode))
		require.ErrorContains(t, theError, contains, msgAndArgs...)
	})
}

// ErrorContains asserts that a function returned an error (i.e. not `nil`)
// and that the error contains the specified substring.
//
// 	actualObj, err := SomeFunction()
// 	assert.ErrorContains(t, err,  expectedErrorSubString)
func (x Assertions[T]) ErrorContains(theError error, contains string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error contains")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("the error", theError).withMode(x.mode), NewParameter("contains", contains).withMode(x.mode))
		assert.ErrorContains(t, theError, contains, msgAndArgs...)
	})
}

// ErrorIs asserts that at least one of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (x Requirements[T]) ErrorIs(err error, target error, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error is")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		require.ErrorIs(t, err, target, msgAndArgs...)
	})
}

// ErrorIs asserts that at least one of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (x Assertions[T]) ErrorIs(err error, target error, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "error is")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		assert.ErrorIs(t, err, target, msgAndArgs...)
	})
}

// Eventually asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
// 	assert.Eventually(t, func() bool { return true; }, time.Second, 10*time.Millisecond)
func (x Requirements[T]) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "eventually")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("condition", condition).withMode(x.mode), NewParameter("wait for", waitFor).withMode(x.mode), NewParameter("tick", tick).withMode(x.mode))
		require.Eventually(t, condition, waitFor, tick, msgAndArgs...)
	})
}

// Eventually asserts that given condition will be met in waitFor time,
// periodically checking target function each tick.
//
// 	assert.Eventually(t, func() bool { return true; }, time.Second, 10*time.Millisecond)
func (x Assertions[T]) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "eventually")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("condition", condition).withMode(x.mode), NewParameter("wait for", waitFor).withMode(x.mode), NewParameter("tick", tick).withMode(x.mode))
		assert.Eventually(t, condition, waitFor, tick, msgAndArgs...)
	})
}

// Exactly asserts that two objects are equal in value and type.
//
// 	assert.Exactly(t, int32(123), int64(123))
func (x Requirements[T]) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "exactly")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.Exactly(t, expected, actual, msgAndArgs...)
	})
}

// Exactly asserts that two objects are equal in value and type.
//
// 	assert.Exactly(t, int32(123), int64(123))
func (x Assertions[T]) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "exactly")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.Exactly(t, expected, actual, msgAndArgs...)
	})
}

// False asserts that the specified value is false.
//
// 	assert.False(t, myBool)
func (x Requirements[T]) False(value bool, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "false")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("value", value).withMode(x.mode))
		require.False(t, value, msgAndArgs...)
	})
}

// False asserts that the specified value is false.
//
// 	assert.False(t, myBool)
func (x Assertions[T]) False(value bool, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "false")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("value", value).withMode(x.mode))
		assert.False(t, value, msgAndArgs...)
	})
}

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (x Requirements[T]) FileExists(path string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "file exists")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		require.FileExists(t, path, msgAndArgs...)
	})
}

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func (x Assertions[T]) FileExists(path string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "file exists")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		assert.FileExists(t, path, msgAndArgs...)
	})
}

// Greater asserts that the first element is greater than the second
//
// 	assert.Greater(t, 2, 1)
// 	assert.Greater(t, float64(2), float64(1))
// 	assert.Greater(t, "b", "a")
func (x Requirements[T]) Greater(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "greater")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		require.Greater(t, e1, e2, msgAndArgs...)
	})
}

// Greater asserts that the first element is greater than the second
//
// 	assert.Greater(t, 2, 1)
// 	assert.Greater(t, float64(2), float64(1))
// 	assert.Greater(t, "b", "a")
func (x Assertions[T]) Greater(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "greater")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		assert.Greater(t, e1, e2, msgAndArgs...)
	})
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
// 	assert.GreaterOrEqual(t, 2, 1)
// 	assert.GreaterOrEqual(t, 2, 2)
// 	assert.GreaterOrEqual(t, "b", "a")
// 	assert.GreaterOrEqual(t, "b", "b")
func (x Requirements[T]) GreaterOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "greater or equal")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		require.GreaterOrEqual(t, e1, e2, msgAndArgs...)
	})
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
// 	assert.GreaterOrEqual(t, 2, 1)
// 	assert.GreaterOrEqual(t, 2, 2)
// 	assert.GreaterOrEqual(t, "b", "a")
// 	assert.GreaterOrEqual(t, "b", "b")
func (x Assertions[T]) GreaterOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "greater or equal")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		assert.GreaterOrEqual(t, e1, e2, msgAndArgs...)
	})
}

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
// 	assert.HTTPBodyContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (x Requirements[T]) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP body contains")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		require.HTTPBodyContains(t, handler, method, url, values, str, msgAndArgs...)
	})
}

// HTTPBodyContains asserts that a specified handler returns a
// body that contains a string.
//
// 	assert.HTTPBodyContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (x Assertions[T]) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP body contains")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		assert.HTTPBodyContains(t, handler, method, url, values, str, msgAndArgs...)
	})
}

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
// 	assert.HTTPBodyNotContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (x Requirements[T]) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP body not contains")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		require.HTTPBodyNotContains(t, handler, method, url, values, str, msgAndArgs...)
	})
}

// HTTPBodyNotContains asserts that a specified handler returns a
// body that does not contain a string.
//
// 	assert.HTTPBodyNotContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")
//
// Returns whether the assertion was successful (true) or not (false).
func (x Assertions[T]) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP body not contains")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		assert.HTTPBodyNotContains(t, handler, method, url, values, str, msgAndArgs...)
	})
}

// HTTPError asserts that a specified handler returns an error status code.
//
// 	assert.HTTPError(t, myHandler, "POST", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func (x Requirements[T]) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP error")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode))
		require.HTTPError(t, handler, method, url, values, msgAndArgs...)
	})
}

// HTTPError asserts that a specified handler returns an error status code.
//
// 	assert.HTTPError(t, myHandler, "POST", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func (x Assertions[T]) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP error")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode))
		assert.HTTPError(t, handler, method, url, values, msgAndArgs...)
	})
}

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
// 	assert.HTTPRedirect(t, myHandler, "GET", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func (x Requirements[T]) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP redirect")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode))
		require.HTTPRedirect(t, handler, method, url, values, msgAndArgs...)
	})
}

// HTTPRedirect asserts that a specified handler returns a redirect status code.
//
// 	assert.HTTPRedirect(t, myHandler, "GET", "/a/b/c", url.Values{"a": []string{"b", "c"}}
//
// Returns whether the assertion was successful (true) or not (false).
func (x Assertions[T]) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP redirect")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode))
		assert.HTTPRedirect(t, handler, method, url, values, msgAndArgs...)
	})
}

// HTTPStatusCode asserts that a specified handler returns a specified status code.
//
// 	assert.HTTPStatusCode(t, myHandler, "GET", "/notImplemented", nil, 501)
//
// Returns whether the assertion was successful (true) or not (false).
func (x Requirements[T]) HTTPStatusCode(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP status code")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode), NewParameter("statuscode", statuscode).withMode(x.mode))
		require.HTTPStatusCode(t, handler, method, url, values, statuscode, msgAndArgs...)
	})
}

// HTTPStatusCode asserts that a specified handler returns a specified status code.
//
// 	assert.HTTPStatusCode(t, myHandler, "GET", "/notImplemented", nil, 501)
//
// Returns whether the assertion was successful (true) or not (false).
func (x Assertions[T]) HTTPStatusCode(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP status code")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode), NewParameter("statuscode", statuscode).withMode(x.mode))
		assert.HTTPStatusCode(t, handler, method, url, values, statuscode, msgAndArgs...)
	})
}

// HTTPSuccess asserts that a specified handler returns a success status code.
//
// 	assert.HTTPSuccess(t, myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func (x Requirements[T]) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP success")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode))
		require.HTTPSuccess(t, handler, method, url, values, msgAndArgs...)
	})
}

// HTTPSuccess asserts that a specified handler returns a success status code.
//
// 	assert.HTTPSuccess(t, myHandler, "POST", "http://www.google.com", nil)
//
// Returns whether the assertion was successful (true) or not (false).
func (x Assertions[T]) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "HTTP success")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("handler", handler).withMode(x.mode), NewParameter("method", method).withMode(x.mode), NewParameter("url", url).withMode(x.mode), NewParameter("values", values).withMode(x.mode))
		assert.HTTPSuccess(t, handler, method, url, values, msgAndArgs...)
	})
}

// Implements asserts that an object is implemented by the specified interface.
//
// 	assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func (x Requirements[T]) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "implements")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("interface object", interfaceObject).withMode(x.mode), NewParameter("object", object).withMode(x.mode))
		require.Implements(t, interfaceObject, object, msgAndArgs...)
	})
}

// Implements asserts that an object is implemented by the specified interface.
//
// 	assert.Implements(t, (*MyInterface)(nil), new(MyObject))
func (x Assertions[T]) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "implements")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("interface object", interfaceObject).withMode(x.mode), NewParameter("object", object).withMode(x.mode))
		assert.Implements(t, interfaceObject, object, msgAndArgs...)
	})
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	assert.InDelta(t, math.Pi, 22/7.0, 0.01)
func (x Requirements[T]) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in delta")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		require.InDelta(t, expected, actual, delta, msgAndArgs...)
	})
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	assert.InDelta(t, math.Pi, 22/7.0, 0.01)
func (x Assertions[T]) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in delta")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		assert.InDelta(t, expected, actual, delta, msgAndArgs...)
	})
}

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (x Requirements[T]) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in delta map values")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		require.InDeltaMapValues(t, expected, actual, delta, msgAndArgs...)
	})
}

// InDeltaMapValues is the same as InDelta, but it compares all values between two maps. Both maps must have exactly the same keys.
func (x Assertions[T]) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in delta map values")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		assert.InDeltaMapValues(t, expected, actual, delta, msgAndArgs...)
	})
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func (x Requirements[T]) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in delta slice")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		require.InDeltaSlice(t, expected, actual, delta, msgAndArgs...)
	})
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func (x Assertions[T]) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in delta slice")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		assert.InDeltaSlice(t, expected, actual, delta, msgAndArgs...)
	})
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func (x Requirements[T]) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in epsilon")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("epsilon", epsilon).withMode(x.mode))
		require.InEpsilon(t, expected, actual, epsilon, msgAndArgs...)
	})
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
func (x Assertions[T]) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in epsilon")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("epsilon", epsilon).withMode(x.mode))
		assert.InEpsilon(t, expected, actual, epsilon, msgAndArgs...)
	})
}

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func (x Requirements[T]) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in epsilon slice")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("epsilon", epsilon).withMode(x.mode))
		require.InEpsilonSlice(t, expected, actual, epsilon, msgAndArgs...)
	})
}

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func (x Assertions[T]) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "in epsilon slice")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("epsilon", epsilon).withMode(x.mode))
		assert.InEpsilonSlice(t, expected, actual, epsilon, msgAndArgs...)
	})
}

// IsDecreasing asserts that the collection is decreasing
//
// 	assert.IsDecreasing(t, []int{2, 1, 0})
// 	assert.IsDecreasing(t, []float{2, 1})
// 	assert.IsDecreasing(t, []string{"b", "a"})
func (x Requirements[T]) IsDecreasing(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is decreasing")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.IsDecreasing(t, object, msgAndArgs...)
	})
}

// IsDecreasing asserts that the collection is decreasing
//
// 	assert.IsDecreasing(t, []int{2, 1, 0})
// 	assert.IsDecreasing(t, []float{2, 1})
// 	assert.IsDecreasing(t, []string{"b", "a"})
func (x Assertions[T]) IsDecreasing(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is decreasing")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.IsDecreasing(t, object, msgAndArgs...)
	})
}

// IsIncreasing asserts that the collection is increasing
//
// 	assert.IsIncreasing(t, []int{1, 2, 3})
// 	assert.IsIncreasing(t, []float{1, 2})
// 	assert.IsIncreasing(t, []string{"a", "b"})
func (x Requirements[T]) IsIncreasing(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is increasing")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.IsIncreasing(t, object, msgAndArgs...)
	})
}

// IsIncreasing asserts that the collection is increasing
//
// 	assert.IsIncreasing(t, []int{1, 2, 3})
// 	assert.IsIncreasing(t, []float{1, 2})
// 	assert.IsIncreasing(t, []string{"a", "b"})
func (x Assertions[T]) IsIncreasing(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is increasing")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.IsIncreasing(t, object, msgAndArgs...)
	})
}

// IsNonDecreasing asserts that the collection is not decreasing
//
// 	assert.IsNonDecreasing(t, []int{1, 1, 2})
// 	assert.IsNonDecreasing(t, []float{1, 2})
// 	assert.IsNonDecreasing(t, []string{"a", "b"})
func (x Requirements[T]) IsNonDecreasing(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is non decreasing")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.IsNonDecreasing(t, object, msgAndArgs...)
	})
}

// IsNonDecreasing asserts that the collection is not decreasing
//
// 	assert.IsNonDecreasing(t, []int{1, 1, 2})
// 	assert.IsNonDecreasing(t, []float{1, 2})
// 	assert.IsNonDecreasing(t, []string{"a", "b"})
func (x Assertions[T]) IsNonDecreasing(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is non decreasing")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.IsNonDecreasing(t, object, msgAndArgs...)
	})
}

// IsNonIncreasing asserts that the collection is not increasing
//
// 	assert.IsNonIncreasing(t, []int{2, 1, 1})
// 	assert.IsNonIncreasing(t, []float{2, 1})
// 	assert.IsNonIncreasing(t, []string{"b", "a"})
func (x Requirements[T]) IsNonIncreasing(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is non increasing")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.IsNonIncreasing(t, object, msgAndArgs...)
	})
}

// IsNonIncreasing asserts that the collection is not increasing
//
// 	assert.IsNonIncreasing(t, []int{2, 1, 1})
// 	assert.IsNonIncreasing(t, []float{2, 1})
// 	assert.IsNonIncreasing(t, []string{"b", "a"})
func (x Assertions[T]) IsNonIncreasing(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is non increasing")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.IsNonIncreasing(t, object, msgAndArgs...)
	})
}

// IsType asserts that the specified objects are of the same type.
//
// 	assert.IsType(t, &MyStruct{}, &MyStruct{})
func (x Requirements[T]) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is type")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected type", expectedType).withMode(x.mode), NewParameter("object", object).withMode(x.mode))
		require.IsType(t, expectedType, object, msgAndArgs...)
	})
}

// IsType asserts that the specified objects are of the same type.
//
// 	assert.IsType(t, &MyStruct{}, &MyStruct{})
func (x Assertions[T]) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "is type")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected type", expectedType).withMode(x.mode), NewParameter("object", object).withMode(x.mode))
		assert.IsType(t, expectedType, object, msgAndArgs...)
	})
}

// JSONEq asserts that two JSON strings are equivalent.
//
// 	assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func (x Requirements[T]) JSONEq(expected string, actual string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "JSON eq")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.JSONEq(t, expected, actual, msgAndArgs...)
	})
}

// JSONEq asserts that two JSON strings are equivalent.
//
// 	assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func (x Assertions[T]) JSONEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "JSON eq")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.JSONEq(t, expected, actual, msgAndArgs...)
	})
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
// 	assert.Len(t, mySlice, 3)
func (x Requirements[T]) Len(object interface{}, length int, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "len")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode), NewParameter("length", length).withMode(x.mode))
		require.Len(t, object, length, msgAndArgs...)
	})
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
// 	assert.Len(t, mySlice, 3)
func (x Assertions[T]) Len(object interface{}, length int, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "len")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode), NewParameter("length", length).withMode(x.mode))
		assert.Len(t, object, length, msgAndArgs...)
	})
}

// Less asserts that the first element is less than the second
//
// 	assert.Less(t, 1, 2)
// 	assert.Less(t, float64(1), float64(2))
// 	assert.Less(t, "a", "b")
func (x Requirements[T]) Less(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "less")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		require.Less(t, e1, e2, msgAndArgs...)
	})
}

// Less asserts that the first element is less than the second
//
// 	assert.Less(t, 1, 2)
// 	assert.Less(t, float64(1), float64(2))
// 	assert.Less(t, "a", "b")
func (x Assertions[T]) Less(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "less")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		assert.Less(t, e1, e2, msgAndArgs...)
	})
}

// LessOrEqual asserts that the first element is less than or equal to the second
//
// 	assert.LessOrEqual(t, 1, 2)
// 	assert.LessOrEqual(t, 2, 2)
// 	assert.LessOrEqual(t, "a", "b")
// 	assert.LessOrEqual(t, "b", "b")
func (x Requirements[T]) LessOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "less or equal")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		require.LessOrEqual(t, e1, e2, msgAndArgs...)
	})
}

// LessOrEqual asserts that the first element is less than or equal to the second
//
// 	assert.LessOrEqual(t, 1, 2)
// 	assert.LessOrEqual(t, 2, 2)
// 	assert.LessOrEqual(t, "a", "b")
// 	assert.LessOrEqual(t, "b", "b")
func (x Assertions[T]) LessOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "less or equal")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e 1", e1).withMode(x.mode), NewParameter("e 2", e2).withMode(x.mode))
		assert.LessOrEqual(t, e1, e2, msgAndArgs...)
	})
}

// Negative asserts that the specified element is negative
//
// 	assert.Negative(t, -1)
// 	assert.Negative(t, -1.23)
func (x Requirements[T]) Negative(e interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "negative")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e", e).withMode(x.mode))
		require.Negative(t, e, msgAndArgs...)
	})
}

// Negative asserts that the specified element is negative
//
// 	assert.Negative(t, -1)
// 	assert.Negative(t, -1.23)
func (x Assertions[T]) Negative(e interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "negative")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e", e).withMode(x.mode))
		assert.Negative(t, e, msgAndArgs...)
	})
}

// Never asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
// 	assert.Never(t, func() bool { return false; }, time.Second, 10*time.Millisecond)
func (x Requirements[T]) Never(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "never")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("condition", condition).withMode(x.mode), NewParameter("wait for", waitFor).withMode(x.mode), NewParameter("tick", tick).withMode(x.mode))
		require.Never(t, condition, waitFor, tick, msgAndArgs...)
	})
}

// Never asserts that the given condition doesn't satisfy in waitFor time,
// periodically checking the target function each tick.
//
// 	assert.Never(t, func() bool { return false; }, time.Second, 10*time.Millisecond)
func (x Assertions[T]) Never(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "never")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("condition", condition).withMode(x.mode), NewParameter("wait for", waitFor).withMode(x.mode), NewParameter("tick", tick).withMode(x.mode))
		assert.Never(t, condition, waitFor, tick, msgAndArgs...)
	})
}

// Nil asserts that the specified object is nil.
//
// 	assert.Nil(t, err)
func (x Requirements[T]) Nil(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "nil")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.Nil(t, object, msgAndArgs...)
	})
}

// Nil asserts that the specified object is nil.
//
// 	assert.Nil(t, err)
func (x Assertions[T]) Nil(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "nil")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.Nil(t, object, msgAndArgs...)
	})
}

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func (x Requirements[T]) NoDirExists(path string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "no dir exists")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		require.NoDirExists(t, path, msgAndArgs...)
	})
}

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func (x Assertions[T]) NoDirExists(path string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "no dir exists")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		assert.NoDirExists(t, path, msgAndArgs...)
	})
}

// NoError asserts that a function returned no error (i.e. `nil`).
//
// 	  actualObj, err := SomeFunction()
// 	  if assert.NoError(t, err) {
// 		   assert.Equal(t, expectedObj, actualObj)
// 	  }
func (x Requirements[T]) NoError(err error, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "no error")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode))
		require.NoError(t, err, msgAndArgs...)
	})
}

// NoError asserts that a function returned no error (i.e. `nil`).
//
// 	  actualObj, err := SomeFunction()
// 	  if assert.NoError(t, err) {
// 		   assert.Equal(t, expectedObj, actualObj)
// 	  }
func (x Assertions[T]) NoError(err error, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "no error")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode))
		assert.NoError(t, err, msgAndArgs...)
	})
}

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func (x Requirements[T]) NoFileExists(path string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "no file exists")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		require.NoFileExists(t, path, msgAndArgs...)
	})
}

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func (x Assertions[T]) NoFileExists(path string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "no file exists")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("path", path).withMode(x.mode))
		assert.NoFileExists(t, path, msgAndArgs...)
	})
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
// 	assert.NotContains(t, "Hello World", "Earth")
// 	assert.NotContains(t, ["Hello", "World"], "Earth")
// 	assert.NotContains(t, {"Hello": "World"}, "Earth")
func (x Requirements[T]) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not contains")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("s", s).withMode(x.mode), NewParameter("contains", contains).withMode(x.mode))
		require.NotContains(t, s, contains, msgAndArgs...)
	})
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
// 	assert.NotContains(t, "Hello World", "Earth")
// 	assert.NotContains(t, ["Hello", "World"], "Earth")
// 	assert.NotContains(t, {"Hello": "World"}, "Earth")
func (x Assertions[T]) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not contains")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("s", s).withMode(x.mode), NewParameter("contains", contains).withMode(x.mode))
		assert.NotContains(t, s, contains, msgAndArgs...)
	})
}

// NotElementsMatch asserts that the specified listA(array, slice...) is NOT equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should not match.
// This is an inverse of ElementsMatch.
//
// assert.NotElementsMatch(t, [1, 1, 2, 3], [1, 1, 2, 3]) -> false
//
// assert.NotElementsMatch(t, [1, 1, 2, 3], [1, 2, 3]) -> true
//
// assert.NotElementsMatch(t, [1, 2, 3], [1, 2, 4]) -> true
func (x Requirements[T]) NotElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not elements match")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list A", listA).withMode(x.mode), NewParameter("list B", listB).withMode(x.mode))
		require.NotElementsMatch(t, listA, listB, msgAndArgs...)
	})
}

// NotElementsMatch asserts that the specified listA(array, slice...) is NOT equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should not match.
// This is an inverse of ElementsMatch.
//
// assert.NotElementsMatch(t, [1, 1, 2, 3], [1, 1, 2, 3]) -> false
//
// assert.NotElementsMatch(t, [1, 1, 2, 3], [1, 2, 3]) -> true
//
// assert.NotElementsMatch(t, [1, 2, 3], [1, 2, 4]) -> true
func (x Assertions[T]) NotElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not elements match")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list A", listA).withMode(x.mode), NewParameter("list B", listB).withMode(x.mode))
		assert.NotElementsMatch(t, listA, listB, msgAndArgs...)
	})
}

// NotEmpty asserts that the specified object is NOT [Empty].
//
// 	if assert.NotEmpty(t, obj) {
// 	  assert.Equal(t, "two", obj[1])
// 	}
func (x Requirements[T]) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not empty")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.NotEmpty(t, object, msgAndArgs...)
	})
}

// NotEmpty asserts that the specified object is NOT [Empty].
//
// 	if assert.NotEmpty(t, obj) {
// 	  assert.Equal(t, "two", obj[1])
// 	}
func (x Assertions[T]) NotEmpty(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not empty")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.NotEmpty(t, object, msgAndArgs...)
	})
}

// NotEqual asserts that the specified values are NOT equal.
//
// 	assert.NotEqual(t, obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (x Requirements[T]) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not equal")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.NotEqual(t, expected, actual, msgAndArgs...)
	})
}

// NotEqual asserts that the specified values are NOT equal.
//
// 	assert.NotEqual(t, obj1, obj2)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func (x Assertions[T]) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not equal")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.NotEqual(t, expected, actual, msgAndArgs...)
	})
}

// NotEqualValues asserts that two objects are not equal even when converted to the same type
//
// 	assert.NotEqualValues(t, obj1, obj2)
func (x Requirements[T]) NotEqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not equal values")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.NotEqualValues(t, expected, actual, msgAndArgs...)
	})
}

// NotEqualValues asserts that two objects are not equal even when converted to the same type
//
// 	assert.NotEqualValues(t, obj1, obj2)
func (x Assertions[T]) NotEqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not equal values")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.NotEqualValues(t, expected, actual, msgAndArgs...)
	})
}

// NotErrorAs asserts that none of the errors in err's chain matches target,
// but if so, sets target to that error value.
func (x Requirements[T]) NotErrorAs(err error, target interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not error as")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		require.NotErrorAs(t, err, target, msgAndArgs...)
	})
}

// NotErrorAs asserts that none of the errors in err's chain matches target,
// but if so, sets target to that error value.
func (x Assertions[T]) NotErrorAs(err error, target interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not error as")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		assert.NotErrorAs(t, err, target, msgAndArgs...)
	})
}

// NotErrorIs asserts that none of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (x Requirements[T]) NotErrorIs(err error, target error, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not error is")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		require.NotErrorIs(t, err, target, msgAndArgs...)
	})
}

// NotErrorIs asserts that none of the errors in err's chain matches target.
// This is a wrapper for errors.Is.
func (x Assertions[T]) NotErrorIs(err error, target error, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not error is")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err", err).withMode(x.mode), NewParameter("target", target).withMode(x.mode))
		assert.NotErrorIs(t, err, target, msgAndArgs...)
	})
}

// NotImplements asserts that an object does not implement the specified interface.
//
// 	assert.NotImplements(t, (*MyInterface)(nil), new(MyObject))
func (x Requirements[T]) NotImplements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not implements")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("interface object", interfaceObject).withMode(x.mode), NewParameter("object", object).withMode(x.mode))
		require.NotImplements(t, interfaceObject, object, msgAndArgs...)
	})
}

// NotImplements asserts that an object does not implement the specified interface.
//
// 	assert.NotImplements(t, (*MyInterface)(nil), new(MyObject))
func (x Assertions[T]) NotImplements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not implements")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("interface object", interfaceObject).withMode(x.mode), NewParameter("object", object).withMode(x.mode))
		assert.NotImplements(t, interfaceObject, object, msgAndArgs...)
	})
}

// NotNil asserts that the specified object is not nil.
//
// 	assert.NotNil(t, err)
func (x Requirements[T]) NotNil(object interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not nil")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		require.NotNil(t, object, msgAndArgs...)
	})
}

// NotNil asserts that the specified object is not nil.
//
// 	assert.NotNil(t, err)
func (x Assertions[T]) NotNil(object interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not nil")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("object", object).withMode(x.mode))
		assert.NotNil(t, object, msgAndArgs...)
	})
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
// 	assert.NotPanics(t, func(){ RemainCalm() })
func (x Requirements[T]) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not panics")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("f", f).withMode(x.mode))
		require.NotPanics(t, f, msgAndArgs...)
	})
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
// 	assert.NotPanics(t, func(){ RemainCalm() })
func (x Assertions[T]) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not panics")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("f", f).withMode(x.mode))
		assert.NotPanics(t, f, msgAndArgs...)
	})
}

// NotRegexp asserts that a specified regexp does not match a string.
//
// 	assert.NotRegexp(t, regexp.MustCompile("starts"), "it's starting")
// 	assert.NotRegexp(t, "^start", "it's not starting")
func (x Requirements[T]) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not regexp")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("rx", rx).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		require.NotRegexp(t, rx, str, msgAndArgs...)
	})
}

// NotRegexp asserts that a specified regexp does not match a string.
//
// 	assert.NotRegexp(t, regexp.MustCompile("starts"), "it's starting")
// 	assert.NotRegexp(t, "^start", "it's not starting")
func (x Assertions[T]) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not regexp")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("rx", rx).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		assert.NotRegexp(t, rx, str, msgAndArgs...)
	})
}

// NotSame asserts that two pointers do not reference the same object.
//
// 	assert.NotSame(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (x Requirements[T]) NotSame(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not same")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.NotSame(t, expected, actual, msgAndArgs...)
	})
}

// NotSame asserts that two pointers do not reference the same object.
//
// 	assert.NotSame(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (x Assertions[T]) NotSame(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not same")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.NotSame(t, expected, actual, msgAndArgs...)
	})
}

// NotSubset asserts that the list (array, slice, or map) does NOT contain all
// elements given in the subset (array, slice, or map).
// Map elements are key-value pairs unless compared with an array or slice where
// only the map key is evaluated.
//
// 	assert.NotSubset(t, [1, 3, 4], [1, 2])
// 	assert.NotSubset(t, {"x": 1, "y": 2}, {"z": 3})
// 	assert.NotSubset(t, [1, 3, 4], {1: "one", 2: "two"})
// 	assert.NotSubset(t, {"x": 1, "y": 2}, ["z"])
func (x Requirements[T]) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not subset")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list", list).withMode(x.mode), NewParameter("subset", subset).withMode(x.mode))
		require.NotSubset(t, list, subset, msgAndArgs...)
	})
}

// NotSubset asserts that the list (array, slice, or map) does NOT contain all
// elements given in the subset (array, slice, or map).
// Map elements are key-value pairs unless compared with an array or slice where
// only the map key is evaluated.
//
// 	assert.NotSubset(t, [1, 3, 4], [1, 2])
// 	assert.NotSubset(t, {"x": 1, "y": 2}, {"z": 3})
// 	assert.NotSubset(t, [1, 3, 4], {1: "one", 2: "two"})
// 	assert.NotSubset(t, {"x": 1, "y": 2}, ["z"])
func (x Assertions[T]) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not subset")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list", list).withMode(x.mode), NewParameter("subset", subset).withMode(x.mode))
		assert.NotSubset(t, list, subset, msgAndArgs...)
	})
}

// NotZero asserts that i is not the zero value for its type.
func (x Requirements[T]) NotZero(i interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not zero")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("i", i).withMode(x.mode))
		require.NotZero(t, i, msgAndArgs...)
	})
}

// NotZero asserts that i is not the zero value for its type.
func (x Assertions[T]) NotZero(i interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "not zero")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("i", i).withMode(x.mode))
		assert.NotZero(t, i, msgAndArgs...)
	})
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
// 	assert.Panics(t, func(){ GoCrazy() })
func (x Requirements[T]) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "panics")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("f", f).withMode(x.mode))
		require.Panics(t, f, msgAndArgs...)
	})
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
// 	assert.Panics(t, func(){ GoCrazy() })
func (x Assertions[T]) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "panics")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("f", f).withMode(x.mode))
		assert.Panics(t, f, msgAndArgs...)
	})
}

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
// 	assert.PanicsWithError(t, "crazy error", func(){ GoCrazy() })
func (x Requirements[T]) PanicsWithError(errString string, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "panics with error")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err string", errString).withMode(x.mode), NewParameter("f", f).withMode(x.mode))
		require.PanicsWithError(t, errString, f, msgAndArgs...)
	})
}

// PanicsWithError asserts that the code inside the specified PanicTestFunc
// panics, and that the recovered panic value is an error that satisfies the
// EqualError comparison.
//
// 	assert.PanicsWithError(t, "crazy error", func(){ GoCrazy() })
func (x Assertions[T]) PanicsWithError(errString string, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "panics with error")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("err string", errString).withMode(x.mode), NewParameter("f", f).withMode(x.mode))
		assert.PanicsWithError(t, errString, f, msgAndArgs...)
	})
}

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
// 	assert.PanicsWithValue(t, "crazy error", func(){ GoCrazy() })
func (x Requirements[T]) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "panics with value")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("f", f).withMode(x.mode))
		require.PanicsWithValue(t, expected, f, msgAndArgs...)
	})
}

// PanicsWithValue asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value equals the expected panic value.
//
// 	assert.PanicsWithValue(t, "crazy error", func(){ GoCrazy() })
func (x Assertions[T]) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "panics with value")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("f", f).withMode(x.mode))
		assert.PanicsWithValue(t, expected, f, msgAndArgs...)
	})
}

// Positive asserts that the specified element is positive
//
// 	assert.Positive(t, 1)
// 	assert.Positive(t, 1.23)
func (x Requirements[T]) Positive(e interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "positive")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e", e).withMode(x.mode))
		require.Positive(t, e, msgAndArgs...)
	})
}

// Positive asserts that the specified element is positive
//
// 	assert.Positive(t, 1)
// 	assert.Positive(t, 1.23)
func (x Assertions[T]) Positive(e interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "positive")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("e", e).withMode(x.mode))
		assert.Positive(t, e, msgAndArgs...)
	})
}

// Regexp asserts that a specified regexp matches a string.
//
// 	assert.Regexp(t, regexp.MustCompile("start"), "it's starting")
// 	assert.Regexp(t, "start...$", "it's not starting")
func (x Requirements[T]) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "regexp")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("rx", rx).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		require.Regexp(t, rx, str, msgAndArgs...)
	})
}

// Regexp asserts that a specified regexp matches a string.
//
// 	assert.Regexp(t, regexp.MustCompile("start"), "it's starting")
// 	assert.Regexp(t, "start...$", "it's not starting")
func (x Assertions[T]) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "regexp")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("rx", rx).withMode(x.mode), NewParameter("str", str).withMode(x.mode))
		assert.Regexp(t, rx, str, msgAndArgs...)
	})
}

// Same asserts that two pointers reference the same object.
//
// 	assert.Same(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (x Requirements[T]) Same(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "same")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.Same(t, expected, actual, msgAndArgs...)
	})
}

// Same asserts that two pointers reference the same object.
//
// 	assert.Same(t, ptr1, ptr2)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func (x Assertions[T]) Same(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "same")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.Same(t, expected, actual, msgAndArgs...)
	})
}

// Subset asserts that the list (array, slice, or map) contains all elements
// given in the subset (array, slice, or map).
// Map elements are key-value pairs unless compared with an array or slice where
// only the map key is evaluated.
//
// 	assert.Subset(t, [1, 2, 3], [1, 2])
// 	assert.Subset(t, {"x": 1, "y": 2}, {"x": 1})
// 	assert.Subset(t, [1, 2, 3], {1: "one", 2: "two"})
// 	assert.Subset(t, {"x": 1, "y": 2}, ["x"])
func (x Requirements[T]) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "subset")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list", list).withMode(x.mode), NewParameter("subset", subset).withMode(x.mode))
		require.Subset(t, list, subset, msgAndArgs...)
	})
}

// Subset asserts that the list (array, slice, or map) contains all elements
// given in the subset (array, slice, or map).
// Map elements are key-value pairs unless compared with an array or slice where
// only the map key is evaluated.
//
// 	assert.Subset(t, [1, 2, 3], [1, 2])
// 	assert.Subset(t, {"x": 1, "y": 2}, {"x": 1})
// 	assert.Subset(t, [1, 2, 3], {1: "one", 2: "two"})
// 	assert.Subset(t, {"x": 1, "y": 2}, ["x"])
func (x Assertions[T]) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "subset")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("list", list).withMode(x.mode), NewParameter("subset", subset).withMode(x.mode))
		assert.Subset(t, list, subset, msgAndArgs...)
	})
}

// True asserts that the specified value is true.
//
// 	assert.True(t, myBool)
func (x Requirements[T]) True(value bool, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "true")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("value", value).withMode(x.mode))
		require.True(t, value, msgAndArgs...)
	})
}

// True asserts that the specified value is true.
//
// 	assert.True(t, myBool)
func (x Assertions[T]) True(value bool, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "true")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("value", value).withMode(x.mode))
		assert.True(t, value, msgAndArgs...)
	})
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
// 	assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second)
func (x Requirements[T]) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "within duration")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		require.WithinDuration(t, expected, actual, delta, msgAndArgs...)
	})
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
// 	assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second)
func (x Assertions[T]) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "within duration")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode), NewParameter("delta", delta).withMode(x.mode))
		assert.WithinDuration(t, expected, actual, delta, msgAndArgs...)
	})
}

// WithinRange asserts that a time is within a time range (inclusive).
//
// 	assert.WithinRange(t, time.Now(), time.Now().Add(-time.Second), time.Now().Add(time.Second))
func (x Requirements[T]) WithinRange(actual time.Time, start time.Time, end time.Time, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "within range")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("actual", actual).withMode(x.mode), NewParameter("start", start).withMode(x.mode), NewParameter("end", end).withMode(x.mode))
		require.WithinRange(t, actual, start, end, msgAndArgs...)
	})
}

// WithinRange asserts that a time is within a time range (inclusive).
//
// 	assert.WithinRange(t, time.Now(), time.Now().Add(-time.Second), time.Now().Add(time.Second))
func (x Assertions[T]) WithinRange(actual time.Time, start time.Time, end time.Time, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "within range")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("actual", actual).withMode(x.mode), NewParameter("start", start).withMode(x.mode), NewParameter("end", end).withMode(x.mode))
		assert.WithinRange(t, actual, start, end, msgAndArgs...)
	})
}

// YAMLEq asserts that two YAML strings are equivalent.
func (x Requirements[T]) YAMLEq(expected string, actual string, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "YAML eq")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		require.YAMLEq(t, expected, actual, msgAndArgs...)
	})
}

// YAMLEq asserts that two YAML strings are equivalent.
func (x Assertions[T]) YAMLEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "YAML eq")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("expected", expected).withMode(x.mode), NewParameter("actual", actual).withMode(x.mode))
		assert.YAMLEq(t, expected, actual, msgAndArgs...)
	})
}

// Zero asserts that i is the zero value for its type.
func (x Requirements[T]) Zero(i interface{}, msgAndArgs ...interface{}) {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "zero")
	Step(x.t, "require: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("i", i).withMode(x.mode))
		require.Zero(t, i, msgAndArgs...)
	})
}

// Zero asserts that i is the zero value for its type.
func (x Assertions[T]) Zero(i interface{}, msgAndArgs ...interface{}) bool {
	x.t.Helper()
	_, file, line, _ := runtime.Caller(1)
	name := cmp.Or(messageFromMsgAndArgs(msgAndArgs...), "zero")
	return testo.Run(x.t, "assert: "+name, func(t T) {
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("%s:%d", file, line)
			}
		})
		t.Parameters(NewParameter("i", i).withMode(x.mode))
		assert.Zero(t, i, msgAndArgs...)
	})
}
