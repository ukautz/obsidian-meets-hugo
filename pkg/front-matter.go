package omh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/gernest/front"
	"gopkg.in/yaml.v2"
)

var (
	ErrNoFrontMatter = errors.New("missing front matter")
)

var frontMatter *front.Matter

// FrontMatter is meta information for markdown documents
type FrontMatter map[string]interface{}

func (fm FrontMatter) Has(key string) bool {
	_, ok := fm[key]
	return ok
}

func (fm FrontMatter) String(key string) string {
	v, ok := fm[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

func (fm FrontMatter) Strings(key string) []string {
	v, ok := fm[key]
	if !ok {
		return nil
	}
	ss, ok := v.([]string)
	if ok {
		return ss
	}
	ii, ok := v.([]interface{})
	if ok {
		ss = make([]string, len(ii))
		for i, vv := range ii {
			s, ok := vv.(string)
			if ok {
				ss[i] = s
			} else {
				ss[i] = fmt.Sprintf("%v", vv)
			}
		}
		return ss
	}
	return nil
}

func ParseFrontMatterMarkdown(content []byte) (FrontMatter, string, error) {
	metaLines := make([]string, 0)
	bodyLines := make([]string, 0)
	state := 0

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if state < 2 && line == "---" {
			state++
			continue
		}
		if state == 1 {
			metaLines = append(metaLines, line)
		} else if state == 2 {
			bodyLines = append(bodyLines, line)
		}
	}
	if len(metaLines) == 0 {
		return nil, "", ErrNoFrontMatter
	}

	meta := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(strings.Join(metaLines, "\n")), &meta)
	if err != nil {
		return nil, "", err
	}

	return FrontMatter(meta), strings.TrimSpace(strings.Join(bodyLines, "\n")), nil
}

func init() {
	frontMatter = front.NewMatter()
	frontMatter.Handle("---", front.YAMLHandler)

}
