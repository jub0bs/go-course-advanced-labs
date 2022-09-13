package twitter

import "testing"

func BenchmarkContainsNoIllegalPattern(b *testing.B) {
	usernames := []string{"", "jub0bs", "abcTwitTerabd"}
	for i := 0; i < b.N; i++ {
		for _, username := range usernames {
			containsNoIllegalPattern(username)
		}
	}
}

func BenchmarkContainsNoIllegalPattern2(b *testing.B) {
	usernames := []string{"", "jub0bs", "abcTwitTerabd"}
	for i := 0; i < b.N; i++ {
		for _, username := range usernames {
			containsNoIllegalPattern2(username)
		}
	}
}
