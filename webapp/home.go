package webapp

import (
	"net/http"
	"text/template"
)

// Index data is used for the home template

type Index struct {
	Title string
}

func IndexHandler(allTemplates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Index{
			Title: "SUSE Console for SAP Applications",
		}
		// you access the cached templates with the defined name, not the filename
		err := allTemplates.ExecuteTemplate(w, "indexPage", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
