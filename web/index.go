package web

import "net/http"

type IndexTemplate struct {
	templateCommon
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplVars :=  IndexTemplate{}
	_, err := initSession(w, r, &tmplVars)

	tmplVars.PageTitle = "Index"

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
