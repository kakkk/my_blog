package misc

func ValPtr[T any](val T) *T {
	return &val
}
