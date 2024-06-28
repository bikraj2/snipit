package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snipit.bikraj.net/internal/models"
	"snipit.bikraj.net/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	User            *models.User
	Form            interface{}
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

func newTemplateCache() (map[string]*template.Template, error) { // Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob() function to get a slice of all filepaths that // match the pattern "./ui/html/pages/*.tmpl". This will essentially gives // us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// Loop through thef page filepaths one-by-one.
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl.html", "html/partials/*.html", page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		// ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		// if err != nil {
		// 	return nil, err
		// }
		// ts, err = ts.ParseFiles(page)
		// fmt.Println(err)
		//   if err != nil {
		// 	return nil, err
		// }
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}
