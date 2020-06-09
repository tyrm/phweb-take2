package web

import "net/http"


type IndexTemplate struct {
	TemplateCommon
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplVars := &IndexTemplate{}

	tmplVars.AlertWarn = &TemplateAlert{
		Text: "stuff",
	}

	t, err := CompileTemplate("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, &tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
