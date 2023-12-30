package bones

import (
	"html/template"
	"net/http"
)

/*
Head struct contains things that should go in to a <head>

Title -  <title>Value</title>
Stylesheet - Slice of strings, each string a "link" to a CSS stylesheet
*/
type Head struct {
	Meta Meta

	Title      string
	Stylesheet []string
}

type Meta struct {
	Charset  string
	Viewport string
}

func DefaultHead() *Head {
	return &Head{
		Meta: Meta{
			Charset: "utf-8",
		},
		Title:      "",
		Stylesheet: []string{},
	}
}

func (h *Head) ToTemplate(w http.ResponseWriter) (*template.Template, error) {
	headTemp, err := template.ParseFiles("./head.tmpl")
	if err != nil {
		return nil, err
	}

	return headTemp, nil
}
