package webapp

import (
	"embed"
	"net/http"
	"text/template"
)

// Index data is used for the home template

type Index struct {
	Title string
}

func IndexHandler(templateFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Index{
			Title: "SUSE Console for SAP Applications",
		}

		tmpl, err := template.ParseFS(templateFS, "templates/home.html.tmpl", "templates/base.html.tmpl")
		// check your err
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
