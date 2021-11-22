package omh

import (
	"bytes"
	"fmt"

	"github.com/gernest/front"
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
	raw, body, err := frontMatter.Parse(bytes.NewReader(content))
	if err != nil {
		return FrontMatter{}, "", err
	}

	return FrontMatter(raw), body, nil
}

func init() {
	frontMatter = front.NewMatter()
	frontMatter.Handle("---", front.YAMLHandler)
}
