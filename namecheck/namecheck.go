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
