package maputil

func Keys[K comparable, V any, M ~map[K]V](m M) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
