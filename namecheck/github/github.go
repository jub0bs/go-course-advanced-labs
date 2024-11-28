package github

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jub0bs/namecheck"
)

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

func (*GitHub) String() string {
	return "GitHub"
}

func (*GitHub) IsValid(username string) bool {
	return isLongEnough(username) &&
		isShortEnough(username) &&
		containsNoIllegalPattern(username) &&
		containsOnlyLegalChars(username) &&
		containsNoIllegalPrefix(username) &&
		containsNoIllegalSuffix(username)
}

func (gh *GitHub) IsAvailable(username string) (bool, error) {
	endpoint := fmt.Sprintf("https://githsdfgsdfgdsfgub.com/%s", url.PathEscape(username))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}
	resp, err := gh.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		return true, nil
	case http.StatusOK:
		return false, nil
	default:
		const tmpl = "namecheck/github: unexpected status code %d"
		return false, fmt.Errorf(tmpl, resp.StatusCode)
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
