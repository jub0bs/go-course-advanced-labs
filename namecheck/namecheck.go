package namecheck

import (
	"context"
	"fmt"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Validator interface {
	IsValid(username string) bool
}

type Availabler interface {
	IsAvailable(ctx context.Context, username string) (bool, error)
}

type Checker interface {
	Validator
	Availabler
	fmt.Stringer
}

type UnknownAvailabilityError struct {
	Platform string
	Username string
	Err      error
}

func (e *UnknownAvailabilityError) Error() string {
	const tmpl = "unknown availability for username %q on platform %s: %v"
	return fmt.Sprintf(tmpl, e.Username, e.Platform, e.Err)
}

func (e *UnknownAvailabilityError) Unwrap() error {
	return e.Err
}
