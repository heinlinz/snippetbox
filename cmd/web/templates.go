package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.fanthom.net/internal/models"
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
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map for template cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		// Parse base file into template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob *on this template set* to add any partial files
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call parse files *on this template set* to add the page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Add the template set to the map, sing the name of the page
		cache[name] = ts
	}

	return cache, nil
}
