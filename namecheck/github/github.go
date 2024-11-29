// Package github allows you to do ...
package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jub0bs/namecheck"
)

// A GitHub is an abstraction for checking validity ....
type GitHub struct {
	Client namecheck.Client
}

const (
	minLen         = 1
	maxLen         = 39
	illegalPrefix  = "-"
	illegalSuffix  = "-"
	illegalPattern = "--"
)

var legalPattern = regexp.MustCompile("^[-0-9A-Za-z]*$")

// String returns a textual representation...
func (*GitHub) String() string {
	return "GitHub"
}

// IsValid...
func (*GitHub) IsValid(username string) bool {
	return isLongEnough(username) &&
		isShortEnough(username) &&
		containsNoIllegalPattern(username) &&
		containsOnlyLegalChars(username) &&
		containsNoIllegalPrefix(username) &&
		containsNoIllegalSuffix(username)
}

var _ namecheck.Checker = (*GitHub)(nil)

// IsAvailable reports the availability of username.
// If the availability cannot be checked, it returns a non-nil error.
func (gh *GitHub) IsAvailable(ctx context.Context, username string) (bool, error) {
	endpoint := fmt.Sprintf("https://github.com/%s", url.PathEscape(username))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	errua := namecheck.UnknownAvailabilityError{
		Username: username,
		Platform: gh.String(),
	}
	if err != nil {
		errua.Err = err
		return false, &errua
	}
	resp, err := gh.Client.Do(req)
	if err != nil {
		errua.Err = err
		return false, &errua
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		return true, nil
	case http.StatusOK:
		return false, nil
	default:
		const tmpl = "namecheck/github: unexpected status code %d"
		errua.Err = fmt.Errorf(tmpl, resp.StatusCode)
		return false, &errua
	}
}

func isLongEnough(username string) bool {
	return utf8.RuneCountInString(username) >= minLen
}

func isShortEnough(username string) bool {
	return utf8.RuneCountInString(username) <= maxLen
}

func containsNoIllegalPattern(username string) bool {
	return !strings.Contains(username, illegalPattern)
}

func containsOnlyLegalChars(username string) bool {
	return legalPattern.MatchString(username)
}

func containsNoIllegalPrefix(username string) bool {
	return !strings.HasPrefix(username, illegalPrefix)
}

func containsNoIllegalSuffix(username string) bool {
	return !strings.HasSuffix(username, illegalSuffix)
}
