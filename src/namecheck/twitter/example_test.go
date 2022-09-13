package twitter_test

import (
	"fmt"

	"github.com/jub0bs/namecheck/twitter"
)

func ExampleTwitter_IsValid() {
	var tw twitter.Twitter
	fmt.Println(tw.IsValid("jub0bs"))
	// Output: true
}
