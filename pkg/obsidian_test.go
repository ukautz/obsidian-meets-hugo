package omh_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	omh "github.com/ukautz/obsidian-meets-hugo/pkg"
)

func TestObsidianNote_HugoFrontMatter(t *testing.T) {
	tests := map[string]struct {
		added map[string]interface{}
		from  omh.ObsidianNote
		to    map[string]interface{}
	}{
		"empty": {
			from: omh.ObsidianNote{},
			to:   map[string]interface{}{"title": ""},
		},
		"added": {
			added: map[string]interface{}{"added": "xxx"},
			from:  omh.ObsidianNote{},
			to:    map[string]interface{}{"title": "", "added": "xxx"},
		},
		"title and matter": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"matter": "one",
					"two":    []string{"matters", "too"},
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":  "the-title",
				"matter": "one",
				"two":    []string{"matters", "too"},
			},
		},
		"date from update": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"date updated": "2021-04-08",
					"date created": "2021-03-06",
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":        "the-title",
				"date":         "2021-04-08T00:00:00Z",
				"date updated": "2021-04-08",
				"date created": "2021-03-06",
			},
		},
		"date from create": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"date created": "2021-03-06",
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":        "the-title",
				"date":         "2021-03-06T00:00:00Z",
				"date created": "2021-03-06",
			},
		},
		"date with hour and minute": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"date created": "2021-03-06 13:14",
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":        "the-title",
				"date":         "2021-03-06T13:14:00Z",
				"date created": "2021-03-06 13:14",
			},
		},
		"date with hour and minute and second": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"date created": "2021-03-06 13:14:15",
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":        "the-title",
				"date":         "2021-03-06T13:14:15Z",
				"date created": "2021-03-06 13:14:15",
			},
		},
		"date in RFC3339": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"date created": "2021-03-06T13:14:15Z",
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":        "the-title",
				"date":         "2021-03-06T13:14:15Z",
				"date created": "2021-03-06T13:14:15Z",
			},
		},
		"date in unssuported": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"date created": "2021-03-06T13",
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title":        "the-title",
				"date created": "2021-03-06T13",
			},
		},
		"drop those aliases": {
			from: omh.ObsidianNote{
				FrontMatter: omh.FrontMatter{
					"aliases": []string{"bla", "blub"},
				},
				Title:   "the-title",
				Content: "whatever",
			},
			to: map[string]interface{}{
				"title": "the-title",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			to := test.from.HugoFrontMatter(test.added)
			assert.Equal(t, test.to, to)
		})
	}
}

func TestLoadObsidianNote(t *testing.T) {
	note, err := omh.LoadObsidianNote(filepath.Join("fixtures", "source", "Some Note.md"))
	require.NoError(t, err)

	assert.Equal(t, omh.ObsidianNote{
		Title:   "Some Note",
		Content: "Link to [[Other Note]] and to [[Note Existing Note]] should all be fine.\n\n![[Something Static.txt]] is also included",
		FrontMatter: omh.FrontMatter{
			"aliases":      []interface{}{"something"},
			"date created": "2021-12-23 11:12:13",
			"date updated": "2021-12-24 11:12:13",
			"tags":         []interface{}{"aaa"},
		},
	}, note)
}

func TestObsidianDirectory_LinkMap(t *testing.T) {
	tests := map[string]struct {
		directory omh.ObsidianDirectory
		linkMap   map[string]string
	}{
		"empty": {
			directory: omh.ObsidianDirectory{},
			linkMap:   map[string]string{},
		},
		"flat": {
			directory: omh.ObsidianDirectory{
				Notes: []omh.ObsidianNote{
					{Title: "Foo Bar"},
					{Title: "Bla Bla"},
				},
			},
			linkMap: map[string]string{

				"Foo Bar": "foo bar/",
				"Bla Bla": "bla bla/",
			},
		},
		"leveled": {
			directory: omh.ObsidianDirectory{
				Notes: []omh.ObsidianNote{
					{Title: "Foo Level 1"},
					{Title: "Bla Level 1"},
				},
				Childs: []omh.ObsidianDirectory{
					{
						Name: "Sub Directory 1",
						Notes: []omh.ObsidianNote{
							{Title: "Foo Level 2a"},
							{Title: "Bla Level 2a"},
						},
						Childs: []omh.ObsidianDirectory{
							{
								Name: "Sub Directory 3",
								Notes: []omh.ObsidianNote{
									{Title: "Foo Level 3"},
									{Title: "Bla Level 3"},
								},
							},
						},
					},
					{
						Name: "Sub Directory 2",
						Notes: []omh.ObsidianNote{
							{Title: "Foo Level 2b"},
							{Title: "Bla Level 2b"},
						},
					},
				},
			},
			linkMap: map[string]string{
				"Bla Level 1":  "bla level 1/",
				"Bla Level 2a": "sub directory 1/bla level 2a/",
				"Bla Level 2b": "sub directory 2/bla level 2b/",
				"Bla Level 3":  "sub directory 1/sub directory 3/bla level 3/",
				"Foo Level 1":  "foo level 1/",
				"Foo Level 2a": "sub directory 1/foo level 2a/",
				"Foo Level 2b": "sub directory 2/foo level 2b/",
				"Foo Level 3":  "sub directory 1/sub directory 3/foo level 3/",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.linkMap, test.directory.LinkMap(strings.ToLower))
		})
	}
}

func TestLoadObsidianDirectory(t *testing.T) {
	directory, err := omh.LoadObsidianDirectory(filepath.Join("fixtures", "source", "Sub Directory"), nil, false)
	require.NoError(t, err)

	assert.Equal(t, "Sub Directory", directory.Name)
	assert.Empty(t, directory.Childs)
	assert.Equal(t, []string{"Circle Thing.svg", "Something Static.txt"}, directory.Files)
	require.Len(t, directory.Notes, 1)
	assert.Equal(t, "Additional Note", directory.Notes[0].Title)
}

func TestLoadObsidianDirectory_Recursive(t *testing.T) {
	directory, err := omh.LoadObsidianDirectory(filepath.Join("fixtures", "source"), nil, true)
	require.NoError(t, err)

	assert.Equal(t, "source", directory.Name)
	assert.Empty(t, directory.Files)

	require.Len(t, directory.Notes, 2)
	assert.Equal(t, "Other Note", directory.Notes[0].Title)
	assert.Equal(t, "Some Note", directory.Notes[1].Title)

	require.Len(t, directory.Childs, 1)
	assert.Equal(t, "Sub Directory", directory.Childs[0].Name)
	assert.Empty(t, directory.Childs[0].Childs)
	assert.Equal(t, []string{"Circle Thing.svg", "Something Static.txt"}, directory.Childs[0].Files)
	require.Len(t, directory.Childs[0].Notes, 1)
	assert.Contains(t, "Additional Note", directory.Childs[0].Notes[0].Title)
}

func TestLoadObsidianDirectory_NotRecursive(t *testing.T) {
	directory, err := omh.LoadObsidianDirectory(filepath.Join("fixtures", "source"), nil, false)
	require.NoError(t, err)

	assert.Equal(t, "source", directory.Name)
	assert.Empty(t, directory.Files)

	require.Len(t, directory.Notes, 2)
	assert.Equal(t, "Other Note", directory.Notes[0].Title)
	assert.Equal(t, "Some Note", directory.Notes[1].Title)

	require.Len(t, directory.Childs, 0)
}

func TestLoadObsidianDirectory_Filtered(t *testing.T) {
	directory, err := omh.LoadObsidianDirectory(filepath.Join("fixtures", "source"), func(on omh.ObsidianNote) bool {
		return on.FrontMatter.Has("more matter")
	}, true)
	require.NoError(t, err)

	assert.Equal(t, "source", directory.Name)
	assert.Empty(t, directory.Files)

	require.Len(t, directory.Notes, 1)
	assert.Equal(t, "Other Note", directory.Notes[0].Title)

	require.Len(t, directory.Childs, 1)
	assert.Equal(t, "Sub Directory", directory.Childs[0].Name)
	require.Len(t, directory.Childs[0].Notes, 1)
	assert.Equal(t, "Additional Note", directory.Childs[0].Notes[0].Title)
}
