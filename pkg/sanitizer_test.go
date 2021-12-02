package omh_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	omh "github.com/ukautz/obsidian-meets-hugo/pkg"
)

func TestSanitize(t *testing.T) {
	tests := map[string]struct {
		from   string
		expect string
	}{
		"empty":    {"", ""},
		"identity": {"a-B-1", "a-B-1"},
		"no space": {"a B 1", "aB1"},
		"no quote": {"a'B'1", "aB1"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := omh.Sanitize(test.from)
			assert.Equal(t, test.expect, v)
		})
	}
}
