package namecheck

import (
	"fmt"
	"net/http"
)

type Validator interface {
	IsValid(username string) bool
}

type Availabler interface {
	IsAvailable(username string) (bool, error)
}

type Checker interface {
	Validator
	Availabler
	fmt.Stringer
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
