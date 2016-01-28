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

	if strings.HasSuffix(err.Error(), "timed out") {
		return true
	}

	if strings.HasSuffix(err.Error(), "i/o timeout") {
		return true
	}

	return false
}
