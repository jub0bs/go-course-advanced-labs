package stub

import (
	"net/http"

	"github.com/jub0bs/namecheck"
)

type clientFunc func(*http.Request) (*http.Response, error)

func (f clientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func ClientWithStatusCode(sc int) namecheck.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		res := http.Response{
			StatusCode: sc,
			Body:       http.NoBody,
		}
		return &res, nil
	}
	return clientFunc(do)
}

func ClientWithError(err error) namecheck.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		return nil, err
	}
	return clientFunc(do)
}
