package service

import (
	"errors"
	"fmt"
	"net/http"
)

var httpMaxCode = 599

type AppError struct {
	StatusCode int
	Status     string
	ErrMessage string
}

func (err AppError) Error() string {
	return fmt.Sprintf("%s: %s", err.Status, err.ErrMessage)
}

func (err AppError) Template() (string, error) {

	if err.StatusCode <= 0 {
		return "", errors.New("illegal status code! is this really an error?")
	}

	if err.StatusCode < http.StatusBadRequest || err.StatusCode > httpMaxCode {
		return "xxx.go.html", nil
	}

	return fmt.Sprintf("%dxx.go.html", int(err.StatusCode/100.0)), nil
}
