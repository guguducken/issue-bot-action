package util

import "errors"

var (
	Notfound       error = errors.New(`404: resource not found`)
	EndLittleStart error = errors.New(`the end time is little to start time`)
	StatusError    error = errors.New(`http status code is invalid, please your settings`)
	TimeoutErr     error = errors.New(`the request is timeout, we will skip this and you can retry`)
	NertworkErr    error = errors.New(`the network is not working now, abort this peocess`)
)
