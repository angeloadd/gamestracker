package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/angeloadd/gamestracker/internal/config"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
)

type Renderer struct {
	cfg           config.Config
	log           *slog.Logger
	TemplateCache map[string]*template.Template
}

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	//Form      *forms.Form
}

var functions = template.FuncMap{}

// NewRenderer sets the config for the template package
func NewRenderer(cfg config.Config, log *slog.Logger) *Renderer {
	return &Renderer{
		cfg: cfg,
		log: log,
	}
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *TemplateData, _ *http.Request) *TemplateData {
	//td.Flash = app.Session.PopString(r.Context(), "flash")
	//td.Warning = app.Session.PopString(r.Context(), "warning")
	//td.Error = app.Session.PopString(r.Context(), "error")
	//td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders a template
func (r *Renderer) Template(w http.ResponseWriter, req *http.Request, tmpl string, td *TemplateData) error {
	var tc map[string]*template.Template

	if r.cfg.IsProduction() {
		tc = r.TemplateCache
	} else {
		tc, _ = r.CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, req)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		r.log.Error("error writing template to browser", err)
		return err
	}

	return nil

}

func (r *Renderer) CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pathToTemplates := path.Join("web", "templates")

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	if len(layouts) == 0 {
		return myCache, fmt.Errorf("no layouts found in %s", pathToTemplates)
	}

	for _, page := range pages {
		ts, err := template.New("base.layout.gohtml").Funcs(functions).ParseFiles(append(layouts, page)...)
		if err != nil {
			return myCache, err
		}

		name := filepath.Base(page)
		r.log.Info(fmt.Sprintf("adding layout %s to template cache", name))
		myCache[name] = ts

	}

	return myCache, nil
}
