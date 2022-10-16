package github_test

import (
	"testing"

	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
)

var _ namecheck.Checker = (*github.GitHub)(nil)

func TestUsernameTooLong(t *testing.T) {
	var gh github.GitHub
	username := "obviously-longer-than-39-chars-skjdhsdkhfkshkfshdkjfhksdjhf"
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf(
			"IsValid(%s) = %t; want %t",
			username,
			got,
			want,
		)
	}
}
