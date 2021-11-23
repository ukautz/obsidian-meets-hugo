package omh_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	omh "github.com/ukautz/obsidian-meets-hugo/pkg"
)

func TestConverter_Run(t *testing.T) {
	output := filepath.Join("fixtures", "dest")
	defer os.RemoveAll(output)

	root, err := omh.LoadObsidianDirectory(filepath.Join("fixtures", "source"), nil, true)
	require.NoError(t, err)

	converter := omh.Converter{
		ConvertName:  strcase.ToKebab,
		ObsidianRoot: root,
		HugoRoot:     output,
		SubPath:      "sub-path",
		FrontMatter: map[string]interface{}{
			"add": "me",
		},
		TagsKey: "alt-tags",
	}
	err = converter.Run()
	require.NoError(t, err)

	expect := filepath.Join("fixtures", "expect")
	assert.Equal(t,
		stripMap(expect, loadDir(t, expect)),
		stripMap(output, loadDir(t, output)),
	)
}

func loadDir(t *testing.T, path string) map[string]string {

	files := make(map[string]string)
	fhs, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatalf("failed read dir %s: %s", path, err)
	}

	for _, fh := range fhs {
		fp := filepath.Join(path, fh.Name())
		if fh.IsDir() {
			sub := loadDir(t, fp)
			for file, content := range sub {
				files[file] = content
			}
		} else {
			raw, err := ioutil.ReadFile(fp)
			if err != nil {
				t.Fatalf("failed to read file %s: %s", fp, err)
			}
			files[fp] = string(raw)
		}

	}

	return files
}

func stripMap(prefix string, from map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range from {
		res[strings.TrimPrefix(k, prefix)] = v
	}
	return res
}
