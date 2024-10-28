package utils

import (
	"context"
	"errors"
	"net/url"
)

func IsTimeout(err error) bool {
	var uErr *url.Error
	switch {
	case errors.Is(err, context.Canceled):
		return true
	case errors.Is(err, context.DeadlineExceeded):
		return true
	case errors.As(err, &uErr) && uErr.Timeout():
		return true
	}

	return false
}
