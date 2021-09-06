package template

import (
	html_template "html/template"
	"os"
	"path/filepath"
	"strings"
	text_template "text/template"

	"github.com/sirupsen/logrus"
)

var parsers = map[string]TemplateBuilder{
	".gohtml": func(filenames ...string) (Template, error) { return html_template.ParseFiles(filenames...) },
	".gotext": func(filenames ...string) (Template, error) { return text_template.ParseFiles(filenames...) },
}

type ProcessPath func(path string, info os.FileInfo) error

func walkDir(dir string, recurse bool, process ProcessPath) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if recurse || path == dir {
				return nil // process root
			}
			return filepath.SkipDir // skip sub-directory
		}
		return process(path, info)
	})
}

type SimpleTemplateManager struct {
	templates map[string]Template
}

func (m *SimpleTemplateManager) Get(name string) (Template, bool) {
	t, found := m.templates[name]
	return t, found
}

func (m *SimpleTemplateManager) Load(name, path string, ParseFiles TemplateBuilder) error {
	if t, err := ParseFiles(path); err != nil {
		return err
	} else {
		m.templates[name] = t
		logrus.Debugf("loaded template[%s] from [%s]\n", name, path)
	}
	return nil
}

func fileNameExt(path string) (name string, ext string) {
	ext = filepath.Ext(path)
	name = strings.TrimSuffix(path, ext)
	return name, ext
}

func (m *SimpleTemplateManager) LoadDir(dir string) error {
	return walkDir(dir, false, func(path string, info os.FileInfo) error {
		name, fileExt := fileNameExt(info.Name())
		parser, exists := parsers[fileExt]
		if !exists {
			return nil
		}
		return m.Load(name, path, parser)
	})
}

func NewSimpleTemplateManager() TemplateManager {
	return &SimpleTemplateManager{templates: make(map[string]Template)}
}
