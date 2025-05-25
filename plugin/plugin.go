package plugin

import (
	"fmt"
	"reflect"

	"testman/internal/reflectutil"
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
	Value any
}

// Pluginer is an interface that plugins may implement to provide
// [Plan], [Hooks] and [Overrides] to the tests.
type Pluginer interface {
	Plugin() Plugin
}

// Plugin specification.
type Plugin struct {
	Plan      Plan
	Hooks     Hooks
	Overrides Overrides
}

// Merge multiple plugins into one.
func Merge(plugins ...Plugin) Plugin {
	return Plugin{
		Plan:      mergePlans(plugins...),
		Hooks:     mergeHooks(plugins...),
		Overrides: mergeOverrides(plugins...),
	}
}

// Collect plugins from the v.
//
// Plugins are stored as (possibly anonymous) fields and implement [Pluginer] interface.
//
// If v itself implements [Pluginer] interface it will
// collect it first and then traverse through its fields recursively.
func Collect(v any) []Plugin {
	return collect(reflect.ValueOf(v))
}

func scan(v any) (Plugin, bool) {
	p, ok := v.(Pluginer)
	if !ok || reflectutil.IsPromotedMethod(reflect.TypeOf(v), "Plugin") {
		return Plugin{}, false
	}

	return p.Plugin(), true
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
