package utils

func MustArg[P any, T any](fn func(P) (T, error), param P) T {
	v, err := fn(param)
	if err != nil {
		panic(err)
	}
	return v
}

func Must[T any](fn func() (T, error)) T {
	v, err := fn()
	if err != nil {
		panic(err)
	}
	return v
}
