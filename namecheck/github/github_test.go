package github_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/stub"
)

func ExampleGitHub_IsValid() {
	var gh github.GitHub
	fmt.Println(gh.IsValid("axel"))
	// Output: true
}

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

func TestIsAvailable200(t *testing.T) {
	gh := github.GitHub{
		Client: stub.ClientWithStatusCode(http.StatusOK),
	}
	avail, err := gh.IsAvailable(context.Background(), "whatever")
	if err != nil || avail {
		t.Errorf("IsAvailable(): got %t, %v; want false, nil", avail, err)
	}
}

func TestIsAvailable404(t *testing.T) {
	gh := github.GitHub{
		Client: stub.ClientWithStatusCode(http.StatusNotFound),
	}
	avail, err := gh.IsAvailable(context.Background(), "whatever")
	if err != nil || !avail {
		t.Errorf("IsAvailable(): got %t, %v; want true, nil", avail, err)
	}
}
