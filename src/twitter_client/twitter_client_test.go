package twitter_client

import (
	"testing"
)

func TestGetUnixTime(t *testing.T) {
	cases := []struct {
		s1  string
		val int
	}{
		{"Sun Feb 15 07:34:13 +0000 2015", 1423985653},
		{"Thu Sep 24 19:37:30 +0000 2009", 1253821050},
	}
	for _, c := range cases {
		got := get_unix_time(c.s1)
		if got != c.val {
			t.Errorf("get_unix_time(%q) == %q, want %q", c.s1, got, c.val)
		}
	}
}

// Most other functions need more of a functional test as they rely on twitter.
// time permitting add them.
