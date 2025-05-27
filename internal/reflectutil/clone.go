package reflectutil

import (
	"reflect"
	"strings"
)

func DeepClone[T any](v T) T {
	cloneMap := make(map[uintptr]reflect.Value)

	return deepClone(reflect.ValueOf(v), cloneMap).Interface().(T)
}

func deepClone(v reflect.Value, cloneMap map[uintptr]reflect.Value) reflect.Value {
	if !v.IsValid() {
		return reflect.Value{}
	}

	switch v.Kind() {
	case reflect.String:
		return reflect.ValueOf(strings.Clone(v.Interface().(string)))

	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return v

	case reflect.Array:
		cloned := reflect.New(reflect.ArrayOf(v.Len(), v.Type().Elem())).Elem()

		for i := range v.Len() {
			cloned.Index(i).Set(deepClone(v.Index(i), cloneMap))
		}

		return cloned

	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		cloned := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := range v.Len() {
			cloned.Index(i).Set(deepClone(v.Index(i), cloneMap))
		}

		return cloned

	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		cloned := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			clonedKey := deepClone(key, cloneMap)
			clonedValue := deepClone(v.MapIndex(key), cloneMap)
			cloned.SetMapIndex(clonedKey, clonedValue)
		}

		return cloned

	case reflect.Struct:
		cloned := reflect.New(v.Type()).Elem()
		for i := range v.NumField() {
			if v.Type().Field(i).PkgPath == "" { // Exported field
				cloned.Field(i).Set(deepClone(v.Field(i), cloneMap))
			}
		}

		return cloned

	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		ptr := v.Pointer()
		if cloned, ok := cloneMap[ptr]; ok {
			return cloned
		}

		clonedElem := deepClone(v.Elem(), cloneMap)
		clonedPtr := reflect.New(v.Type().Elem())
		clonedPtr.Elem().Set(clonedElem)
		cloneMap[ptr] = clonedPtr

		return clonedPtr

	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}

		return deepClone(v.Elem(), cloneMap)

	default:
		panic("unsupported kind: " + v.Kind().String())
	}
}
