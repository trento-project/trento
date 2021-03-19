package webapp

import (
	"fmt"
	"html/template"
	"net/http"
)

// Index data is used for the home template

type Index struct {
	Title     string
	Copyright string
}

func IndexHandler(templates map[string]*template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Index{
			Title:     "SUSE Console for SAP Applications",
			Copyright: "2019-2021 SUSE, all rights reserved.",
		}

		tmpl, ok := templates["home.html.tmpl"]

		//		err = tmpl.ExecuteTemplate(w, "base", data)
		if !ok {
			http.Error(w, fmt.Sprintf("The template HOME does not exist"),
				http.StatusInternalServerError)
			return
		}

		err := tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
