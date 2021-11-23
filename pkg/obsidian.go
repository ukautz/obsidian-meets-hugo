package omh

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/gernest/front"
	log "github.com/sirupsen/logrus"
)

// ObsidianFilter includes or excludes a note
type ObsidianFilter func(ObsidianNote) bool

// ObsidianNote is a single note in Obsidian
type ObsidianNote struct {
	FrontMatter
	Title     string
	Content   string
	Directory *ObsidianDirectory
}

// HugoFrontMatter returns an updated front-matter metadata, suitable for Hugo pages
func (note ObsidianNote) HugoFrontMatter(added map[string]interface{}) map[string]interface{} {
	hugo := make(map[string]interface{})
	for k, v := range note.FrontMatter {
		hugo[k] = v
	}
	for k, v := range added {
		hugo[k] = v
	}

	// must have title
	hugo["title"] = note.Title

	// if date exists, use that
	if note.Has("date updated") {
		hugo["date"] = hugo["date updated"]
	} else if note.Has("date created") {
		hugo["date"] = hugo["date created"]
	}

	// cleanup
	delete(hugo, "aliases")
	return hugo
}

// LoadObsidianNote loads an Obsidian note from disk at given path
func LoadObsidianNote(path string) (ObsidianNote, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return ObsidianNote{}, err
	}

	matter, content, err := ParseFrontMatterMarkdown(raw)
	if err != nil {
		return ObsidianNote{}, err
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return ObsidianNote{
		FrontMatter: matter,
		Title:       title,
		Content:     content,
	}, nil
}

// ObsidianDirectory is a directory within an Obsidian Vault
type ObsidianDirectory struct {
	Name   string
	Path   string
	Parent *ObsidianDirectory
	Childs []ObsidianDirectory
	Notes  []ObsidianNote
	Files  []string
}

func (directory ObsidianDirectory) Empty() bool {
	return len(directory.Childs) == 0 && len(directory.Files) == 0 && len(directory.Notes) == 0
}

// LinkMap is the map of Obsidian internal links to Hugo compatible web links ({"Internal Name": "directory/internal-name"}).
// Note that the Obsidian structure is flat!
func (directory ObsidianDirectory) LinkMap(convert ConvertName) map[string]string {
	to := make(map[string]string)
	directory.linkMap(convert, to, "")
	return to
}

func (directory ObsidianDirectory) linkMap(convert ConvertName, to map[string]string, prefix string) {
	for _, note := range directory.Notes {
		target := path.Join(prefix, convert(note.Title))
		if _, ok := to[note.Title]; ok && to[note.Title] != target {
			log.WithFields(log.Fields{
				"title":   note.Title,
				"target1": to[note.Title],
				"target2": target,
			}).Warn("duplicate link found (same Obsidian note in different directories?)")
		}
		to[note.Title] = target
	}
	for _, sub := range directory.Childs {
		sub.linkMap(convert, to, path.Join(prefix, convert(sub.Name)))
	}
}

// LoadObsidianDirectory reads all notes and sub-directories within a directory in an Obsidian vault
func LoadObsidianDirectory(path string, filter ObsidianFilter, recurse bool) (root ObsidianDirectory, err error) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	root.Path = path
	root.Name = filepath.Base(path)
	root.Childs = make([]ObsidianDirectory, 0)
	root.Files = make([]string, 0)
	root.Notes = make([]ObsidianNote, 0)
	for _, fi := range fis {

		// ignore hidden
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		p := filepath.Join(path, fi.Name())

		// recurse directories
		if fi.IsDir() {
			if !recurse {
				continue
			}
			log.WithField("directory", p).Debug("traverse sub-directory")
			sub, err := LoadObsidianDirectory(p, filter, true)
			if err != nil {
				return ObsidianDirectory{}, err
			} else if sub.Empty() {
				continue
			}

			sub.Parent = &root
			root.Childs = append(root.Childs, sub)

			// handle markdown files
		} else if filepath.Ext(p) == ".md" {
			log.WithField("file", p).Debug("load markdown file")

			note, err := LoadObsidianNote(p)
			if err != nil {

				// ignore markdown files that lack front-matter
				if errors.Is(err, front.ErrUnknownDelim) {
					log.WithFields(log.Fields{"file": p}).Warn("ignore file with missing front matter")
					continue
				}
				return ObsidianDirectory{}, err
			}

			if filter != nil && !filter(note) {
				log.WithField("note", note.Title).Info("note filtered out")
				continue
			}

			note.Directory = &root
			root.Notes = append(root.Notes, note)

			// handle other (static) files
		} else {
			log.WithField("file", p).Debug("add static file")
			root.Files = append(root.Files, fi.Name())

		}
	}

	return
}
