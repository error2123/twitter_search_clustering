package tweet_similarity

import (
	"testing"
)

func TestCosine(t *testing.T) {
	cases := []struct {
		s1  string
		s2  string
		val float64
	}{
		{"We are same", "we are same", 1.0000000000000002},
		{"We are different", "for sure??", 0.0001},
	}
	for _, c := range cases {
		got := Cosine(c.s1, c.s2)
		if got != c.val {
			t.Errorf("Cosine(%q, %q) == %q, want %q", c.s1, c.s2, got, c.val)
		}
	}
}
