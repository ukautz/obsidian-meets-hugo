package omh_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	omh "github.com/ukautz/obsidian-meets-hugo/pkg"
)

func TestFrontMatter_Has(t *testing.T) {
	m := omh.FrontMatter{"a": 1, "b": "B", "bb": 3}

	assert.True(t, m.Has("a"))
	assert.True(t, m.Has("b"))
	assert.True(t, m.Has("bb"))
	assert.False(t, m.Has("c"))
	assert.False(t, m.Has("B"))
	assert.False(t, m.Has("bbb"))
}

func TestFrontMatter_String(t *testing.T) {
	m := omh.FrontMatter{"a": 1, "b": "B", "bb": []string{"x", "y"}}

	assert.Equal(t, "1", m.String("a"))
	assert.Equal(t, "B", m.String("b"))
	assert.Equal(t, "[x y]", m.String("bb"))
	assert.Equal(t, "", m.String("c"))
	assert.Equal(t, "", m.String("d"))
}

func TestFrontMatter_Strings(t *testing.T) {
	m := omh.FrontMatter{"a": 1, "b": "B", "bb": []string{"x", "y"}, "cc": []interface{}{"y", "z"}}

	assert.Nil(t, m.Strings("a"))
	assert.Nil(t, m.Strings("b"))
	assert.Equal(t, []string{"x", "y"}, m.Strings("bb"))
	assert.Equal(t, []string{"y", "z"}, m.Strings("cc"))
}

func TestParseFrontMatterMarkdown(t *testing.T) {
	fm, body, err := omh.ParseFrontMatterMarkdown([]byte(`---
foo: 1
bar: bla
baz:
  - one
  - two
  - 3
---

and the body and stuff`))

	require.Nil(t, err)
	assert.Equal(t, omh.FrontMatter{
		"foo": 1,
		"bar": "bla",
		"baz": []interface{}{"one", "two", 3},
	}, fm)
	assert.Equal(t, "and the body and stuff", body)
}

func TestParseFrontMatterMarkdown_FrontMatterNotCode(t *testing.T) {
	rawMeta := strings.Join([]string{
		"---",
		"foo: 1",
		"---",
	}, "\n")
	rawBody := strings.Join([]string{
		"Must not be confused:",
		"",
		"```",
		"---",
		"valid: code block",
		"foo: 2",
		"---",
		"```",
		"",
		"Fin",
	}, "\n")

	fm, body, err := omh.ParseFrontMatterMarkdown([]byte(rawMeta + "\n\n" + rawBody))

	require.Nil(t, err)
	assert.Equal(t, omh.FrontMatter{
		"foo": 1,
	}, fm)
	assert.Equal(t, rawBody, body)
}
