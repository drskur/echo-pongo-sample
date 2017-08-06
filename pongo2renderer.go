package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

// PongoRenderer implements custom pongo2 rendering engine for labstack's echo framework
type PongoRenderer struct {
	dirs              []string
	templates         *pongo2.TemplateSet
	contextProcessors []ContextProcessorFunc
}

// PongoRendererOption signature
type PongoRendererOption func(*PongoRenderer)

// ContextProcessorFunc signature
type ContextProcessorFunc func(echoCtx echo.Context, pongoCtx pongo2.Context)

// NewPongoRenderer creates a new renderer
func NewPongoRenderer(options ...PongoRendererOption) *PongoRenderer {
	r := &PongoRenderer{}
	r.templates = pongo2.NewSet("templates", r)

	for _, opt := range options {
		opt(r)
	}

	return r
}

// Debug configures template debugging
func Debug(isDebug bool) PongoRendererOption {
	return func(r *PongoRenderer) {
		r.templates.Debug = isDebug
	}
}

// SetGlobals configures global variables passed to template engine
func SetGlobals(globals map[string]interface{}) PongoRendererOption {
	return func(r *PongoRenderer) {
		for k, v := range globals {
			r.templates.Globals[k] = v
		}
	}
}

// RegisterTag registers a custom tag
func RegisterTag(tagname string, tagfunc pongo2.TagParser) PongoRendererOption {
	return func(r *PongoRenderer) {
		pongo2.RegisterTag(tagname, tagfunc)
	}
}

// UseContextProcessor adds context processor to the pipeline
func (r *PongoRenderer) UseContextProcessor(processor ContextProcessorFunc) {
	r.contextProcessors = append(r.contextProcessors, processor)
}

// Abs returns absolute path to file requested
func (r *PongoRenderer) Abs(base, name string) string {
	if filepath.IsAbs(name) {
		return name
	}

	for _, dir := range r.dirs {
		fullpath := filepath.Join(dir, name)
		_, err := os.Stat(fullpath)
		if err == nil {
			return fullpath
		}
	}

	return filepath.Join(filepath.Dir(base), name)
}

// Get reads the path's content from your local filesystem.
func (r *PongoRenderer) Get(path string) (io.Reader, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

// AddDirectory adds a directory to the list of directories searched for templates
func (r *PongoRenderer) AddDirectory(dir string) {
	r.dirs = append(r.dirs, dir)
}

// Render renders the view
func (r *PongoRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := r.templates.FromCache(name)
	if err != nil {
		return err
	}
	d, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("Incorrect data format. Should be map[string]interface{}")
	}

	// run context processors
	for _, processor := range r.contextProcessors {
		processor(c, d)
	}

	return tmpl.ExecuteWriter(pongo2.Context(d), w)
}
