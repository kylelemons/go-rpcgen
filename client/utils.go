package client

import (
	"strings"
)

// Given an error, makes a crappy attempt to determine if the error is a
// timeout error.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	return strings.HasSuffix(err.Error(), "connection timed out")
}
