package utils

func NotNil[T any](v *T) T {
	if v == nil {
		var def T
		return def
	}

	return *v
}
