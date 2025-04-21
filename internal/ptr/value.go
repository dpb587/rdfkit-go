package ptr

func Value[T any](v T) *T {
	return &v
}
