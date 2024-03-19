package BookQuest

import (
	"html/template"
)

type Template struct {
	Templates *template.Template
}

type FormData struct {
	Value string
}

type Profile struct {
	Id             string
	Name           string
	LinksAdded     int
	LinksClicked   int
	LinksFavoirted int
}

type TemplateData struct {
	Form    FormData
	Results []Page
}

func NewTemplateData() TemplateData {
	return TemplateData{
		Form: FormData{
			Value: "",
		},
		Results: []Page{},
	}
}

func NewTemplateRenderer(paths ...string) template.Template {
	tmpl := template.Template{}
	for i := range paths {
		template.Must(tmpl.ParseGlob(paths[i]))
	}
	return tmpl
}
