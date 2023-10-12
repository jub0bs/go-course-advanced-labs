package stub

import (
	"io"
	"net/http"
	"strings"

	"github.com/jub0bs/namecheck"
)

type clientFunc func(*http.Request) (*http.Response, error)

func (f clientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func ClientWithError(err error) namecheck.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		return nil, err
	}
	return clientFunc(do)
}

func ClientWithStatusCodeAndBody(sc int, body string) namecheck.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		res := http.Response{
			StatusCode: sc,
			Body:       io.NopCloser(strings.NewReader(body)),
		}
		return &res, nil
	}
	return clientFunc(do)
}

func ClientWithStatusCode(sc int) namecheck.Client {
	return ClientWithStatusCodeAndBody(sc, "")
}
