// Package plugin provides plugin primitives for using plugins in testo.
package plugin

import (
	"fmt"
	"reflect"

	"github.com/metafates/testo/internal/reflectutil"
)

// Option is used to configure plugin upon creation.
//
// All user-supplied options are passed to the New method for each plugin.
// It is a plugin responsibility to check if the given option corresponds to it.
// One way to check it is with type assertion:
//
//	var opt Option
//	o, ok := opt.(MyOption)
type Option struct {
	// Value of this option.
	Value any

	// Propagate states whether this option
	// should be passed automatically between all subtests.
	Propagate bool
}

// Plugin is an interface that plugins implement to provide
// [Plan], [Hooks] and [Overrides] to the tests.
type Plugin interface {
	Plugin() Spec
}

// Spec specification.
type Spec struct {
	Plan      Plan
	Hooks     Hooks
	Overrides Overrides
}

// MergeSpecs multiple plugin specs into one.
func MergeSpecs(plugins ...Spec) Spec {
	return Spec{
		Plan:      mergePlans(plugins...),
		Hooks:     mergeHooks(plugins...),
		Overrides: mergeOverrides(plugins...),
	}
}

// Collect plugins from the v.
//
// Plugins are stored as (possibly anonymous) fields and implement [Plugin] interface.
//
// If v itself implements [Plugin] interface it will
// collect it first and then traverse through its fields recursively.
func Collect(v any) []Plugin {
	return collect(reflect.ValueOf(v))
}

func scan(v any) (Plugin, bool) {
	p, ok := v.(Plugin)
	if !ok || reflectutil.IsPromotedMethod(reflect.TypeOf(v), "Plugin") {
		return nil, false
	}

	return p, true
}

func collect(value reflect.Value) []Plugin {
	if value.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("expected value kind to be a pointer, got %s", value.Type()))
	}

	if !value.CanInterface() {
		return nil
	}

	var plugins []Plugin

	p, ok := scan(value.Interface())
	if ok {
		plugins = append(plugins, p)
	}

	value = reflectutil.Elem(value)

	if value.Kind() != reflect.Struct {
		return plugins
	}

	for i := range value.NumField() {
		valueField := value.Field(i)

		if valueField.Kind() == reflect.Pointer {
			plugins = append(plugins, collect(valueField)...)
		} else {
			plugins = append(plugins, collect(valueField.Addr())...)
		}
	}

	return plugins
}
