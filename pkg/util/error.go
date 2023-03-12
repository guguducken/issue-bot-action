package util

import "errors"

var (
	Notfound       error = errors.New(`404: resource not found`)
	EndLittleStart error = errors.New(`the end time is little to start time`)
	StatusError    error = errors.New(`http status code is invalid, please your settings`)
)
