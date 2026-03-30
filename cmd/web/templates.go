package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.fanthom.net/internal/models"
	"snippetbox.fanthom.net/ui"
)

// Define a templateData type to act as a holding structure for
// any dynamic data that we want to pass into the HTML templatess
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map for template cache
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// // Call ParseGlob *on this template set* to add any partial files
		// ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		// if err != nil {
		// 	return nil, err
		// }
		//
		// // Call parse files *on this template set* to add the page
		// ts, err = ts.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }
		// Add the template set to the map, sing the name of the page

		cache[name] = ts
	}

	return cache, nil
}
