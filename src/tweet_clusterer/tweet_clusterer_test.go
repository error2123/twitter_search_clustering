package tweet_clusterer

import (
	"testing"
)

func TestInsert(t *testing.T) {
	cases := []struct {
		in   []string
		idx  int
		val  string
		want []string
	}{{[]string{}, 0, "test", []string{"test"}}}
	for _, c := range cases {
		got := insert(c.in, c.idx, c.val)
		if got[0] != "test" {
			t.Errorf("insert(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestStringConcat(t *testing.T) {
	cases := []struct {
		s1   string
		s2   string
		want string
	}{{"yes?", "no", "yes?no"}}
	for _, c := range cases {
		got := string_concat(c.s1, c.s2)
		if got != c.want {
			t.Errorf("string_concat(%q, %q) == %q, want %q", c.s1, c.s2, got, c.want)
		}
	}
}

func TestRemoveStopwords(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{"a and an", " "}}
	for _, c := range cases {
		got := remove_stopwords(c.in)
		if got != c.want {
			t.Errorf("remove_stopwords(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCleanString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{"/:", " "}}
	for _, c := range cases {
		got := clean_string(c.in)
		if got != c.want {
			t.Errorf("clean_string(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

