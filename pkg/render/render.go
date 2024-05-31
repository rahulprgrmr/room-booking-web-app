package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rahulprgrmr/room-booking-web-app/pkg/config"
	"github.com/rahulprgrmr/room-booking-web-app/pkg/models"
)

func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}

var app *config.AppConfig

// New Templates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	// get the template cache from the app config
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
	// create a template cache
	// templateCache, err := CreateTemplateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatalf("couldn't find the template: %s", tmpl)
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	err := t.Execute(buf, templateData)
	if err != nil {
		log.Println(err)
		return
	}
	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		if name == "" {
			continue
		}
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
	
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

// Below function is a simple approach for template caching

var templateCacheSimple = make(map[string]*template.Template)
func RenderTemplateSimple(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache
	_, inMap := templateCacheSimple[t]
	if !inMap {
		// need to create the template
		log.Println("creating template and adding to cache")
		err = createSimpleTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in the cache
		log.Println("using cached template")
	}

	tmpl = templateCacheSimple[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createSimpleTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache
	templateCacheSimple[t] = tmpl

	return nil
}