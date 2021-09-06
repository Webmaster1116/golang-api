package template

import (
	"io"
)

type Template interface {
	Execute(wr io.Writer, data interface{}) error
}
type TemplateBuilder func(filenames ...string) (Template, error)

type TemplateManager interface {
	// get template by name
	Get(name string) (Template, bool)
	// load go template
	Load(name, path string, parser TemplateBuilder) error
	LoadDir(path string) error
}

func Impl() TemplateManager {
	return NewSimpleTemplateManager()
}
