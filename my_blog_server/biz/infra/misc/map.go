package misc

func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k, _ := range m {
		key := k
		result = append(result, key)
	}
	return result
}

func MapValues[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		val := v
		result = append(result, val)
	}
	return result
}
