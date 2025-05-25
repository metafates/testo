package plugin

import (
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
	var plugins []Plugin

	rootPlugin, ok := scanPlugin(v)
	if ok {
		plugins = append(plugins, rootPlugin)
	}

	rv := reflectutil.Elem(reflect.ValueOf(v))

	if rv.Kind() == reflect.Struct {
		for i := range rv.NumField() {
			field := rv.Field(i)

			if field.IsValid() && rv.Type().Field(i).IsExported() {
				plugins = append(plugins, Collect(field.Interface())...)
			}
		}
	}

	return plugins
}

func scanPlugin(v any) (Plugin, bool) {
	p, ok := v.(interface{ Plugin() Plugin })
	if !ok || reflectutil.IsPromotedMethod(reflect.TypeOf(v), "Plugin") {
		return Plugin{}, false
	}

	return p.Plugin(), true
}
