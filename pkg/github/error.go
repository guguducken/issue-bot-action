package github

import "errors"

var (
	NotModify    error = errors.New(`GitHub: Not modified`)
	InvalidToken error = errors.New(`GitHub: Validation failed, or the endpoint has been spammed.`)

	Forbidden error = errors.New(`GitHub: Forbidden`)
)
