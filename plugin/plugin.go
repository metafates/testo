package plugin

import (
	"reflect"

	"testman/internal/reflectutil"
)

// Option is used to configure plugin upon creation.
//
// All user-supplied options are passed to the New method for each plugin.
// It is a plugin responsibility to check if the given option corresponds to it.
type Option struct {
	Value any
}

type Plugin struct {
	Plan      Plan
	Hooks     Hooks
	Overrides Overrides
}

func Merge(plugins ...Plugin) Plugin {
	return Plugin{
		Plan:      mergePlans(plugins...),
		Hooks:     mergeHooks(plugins...),
		Overrides: mergeOverrides(plugins...),
	}
}

func Collect(v any) []Plugin {
	return collectPlugins(v, make(map[uintptr]struct{}))
}

func collectPlugins(v any, visited map[uintptr]struct{}) []Plugin {
	var plugins []Plugin

	rootPlugin, ok := scanPlugin(v, visited)
	if ok {
		plugins = append(plugins, rootPlugin)
	}

	rv := reflectutil.Elem(reflect.ValueOf(v))

	if rv.Kind() == reflect.Struct {
		for i := range rv.NumField() {
			field := rv.Field(i)

			if field.IsValid() && rv.Type().Field(i).IsExported() {
				plugins = append(plugins, collectPlugins(field.Interface(), visited)...)
			}
		}
	}

	return plugins
}

func scanPlugin(v any, visited map[uintptr]struct{}) (Plugin, bool) {
	// We could (and will) access the same method twice if it was promoted
	// from an embed type, which will result calling it twice later.
	//
	// Therefore we maintain visited pointers set so that we won't collect the same method twice.

	if v, ok := v.(interface{ Plugin() Plugin }); ok {
		ptr := reflect.ValueOf(v.Plugin).Pointer()

		if _, ok := visited[ptr]; !ok {
			visited[ptr] = struct{}{}

			return v.Plugin(), true
		}
	}

	return Plugin{}, false
}
