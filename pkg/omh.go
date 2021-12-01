// Package omh implements converter tooling from Obsidian to Hugo
package omh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ConvertName makes a name link-suitable
type ConvertName func(name string) (link string)

// Converter transforms all notes from Obsidian (Vault) directory into pages in Hugo, with rewritten internal links
type Converter struct {
	ConvertName

	// ObsidianRoot is the root of the Obsidian Vault (or a sub-directory thereof)
	ObsidianRoot ObsidianDirectory

	// HugoRoot is the root of the Hugo setup, which contains `content` and `static` folders
	HugoRoot string

	// SubPath defaults to `obsidian` and is the sub-path that will be used under `content` and `static`
	SubPath string

	// FrontMatter is additional front-matter added to each document
	FrontMatter map[string]interface{}

	// TagsKey is name of the key in front-matter that should contain tags (or unset, in case not changed)
	TagsKey string

	linkMap map[string]string
}

func (c *Converter) init() {
	c.linkMap = c.ObsidianRoot.LinkMap(c.ConvertName)
}

// Run transforms and writes all Obsidian root found Markdown files into Hugo suitable Markdown files as well as copies all used static
func (c *Converter) Run() (err error) {
	c.init()

	err = c.processFiles(c.ObsidianRoot, filepath.Join(c.HugoRoot, "static", c.SubPath))
	if err != nil {
		return
	}

	err = c.processNotes(c.ObsidianRoot, filepath.Join(c.HugoRoot, "content", c.SubPath))

	return
}

func (c Converter) processFiles(obsidianDir ObsidianDirectory, hugoDir string) error {
	err := os.MkdirAll(hugoDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// move all files and make them link-able
	startWord := filepath.Join("static", c.SubPath)
	startLen := len(startWord) + 1
	for _, file := range obsidianDir.Files {
		src := filepath.Join(obsidianDir.Path, file)
		ext := filepath.Ext(file)
		dst := filepath.Join(hugoDir, c.ConvertName(strings.TrimSuffix(file, ext))+ext)
		if err = c.copyFile(src, dst); err != nil {
			return err
		}

		// add to link map, so will be replaced later on
		idx := strings.Index(dst, startWord)
		rel := dst[idx+startLen:]
		c.linkMap[file] = rel
	}

	// recurse
	for _, obsidianSubDir := range obsidianDir.Childs {
		if len(obsidianSubDir.Childs) == 0 && len(obsidianSubDir.Files) == 0 {
			continue
		}

		hugoSubPath := filepath.Join(hugoDir, c.ConvertName(obsidianSubDir.Name))
		if err = c.processFiles(obsidianSubDir, hugoSubPath); err != nil {
			return err
		}
	}

	return nil
}

func (c ConvertName) copyFile(from, to string) error {
	src, err := os.OpenFile(from, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (c Converter) processNotes(obsidianDir ObsidianDirectory, hugoDir string) error {
	err := os.MkdirAll(hugoDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// move all notes
	for _, note := range obsidianDir.Notes {
		hugoContent, err := c.convertNote(note)
		if err != nil {
			return fmt.Errorf("failed to convert %s: %w", note.Title, err)
		}

		hugoPath := filepath.Join(hugoDir, c.ConvertName(note.Title)) + ".md"
		err = ioutil.WriteFile(hugoPath, hugoContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", hugoPath, err)
		}
	}

	// recurse
	for _, obsidianSubDir := range obsidianDir.Childs {
		if len(obsidianSubDir.Childs) == 0 && len(obsidianSubDir.Notes) == 0 {
			continue
		}

		hugoSubPath := filepath.Join(hugoDir, c.ConvertName(obsidianSubDir.Name))
		if err := c.processNotes(obsidianSubDir, hugoSubPath); err != nil {
			return err
		}
	}

	return nil
}

var (
	obsidianLink = regexp.MustCompile(`\[\[.+?\]\]`)
)

func (c Converter) convertNote(note ObsidianNote) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// write front matter
	buf.WriteString("---\n")
	matter := note.HugoFrontMatter(c.FrontMatter)
	if c.TagsKey != "" && c.TagsKey != "tags" {
		tags, ok := matter["tags"]
		if ok {
			matter[c.TagsKey] = tags
			delete(matter, "tags")
		}
	}

	frontMatter, err := yaml.Marshal(matter)
	if err != nil {
		return nil, err
	}
	buf.Write(frontMatter)
	buf.WriteString("---\n\n\n")

	// replace internal links in content with "regular" links
	content := obsidianLink.ReplaceAllStringFunc(note.Content, func(s string) string {
		s = strings.TrimPrefix(s, "[[")
		s = strings.TrimSuffix(s, "]]")

		link, title := s, s
		if i := strings.Index(s, "|"); i > -1 {
			link, title = s[0:i], s[i+1:]
		}

		target, ok := c.linkMap[link]
		if !ok {
			log.WithFields(log.Fields{
				"link-title":  title,
				"link-target": link,
				"note":        note.Title,
			}).Warn("missing target for note")
			return title
		}

		return fmt.Sprintf("[%s](/%s/%s)", title, c.SubPath, target)
	})
	buf.WriteString(content)

	return buf.Bytes(), nil

}
