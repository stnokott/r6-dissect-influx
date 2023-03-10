package utils

import (
	"log"

	"github.com/ncruces/zenity"
)

func ErrDialog(err error) {
	if errInner := zenity.Error(err.Error(), zenity.Title("Error"), zenity.ErrorIcon, zenity.Modal()); errInner != nil {
		log.Println("error while displaying error dialog:", errInner)
		log.Println("original error:", err)
	}
}

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
