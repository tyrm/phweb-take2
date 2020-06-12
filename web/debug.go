package web

import "net/http"

func HandleDebug(w http.ResponseWriter, r *http.Request) {
	tmplVars :=  IndexTemplate{}
	_, err := initSession(w, r, &tmplVars)

	tmplVars.PageTitle = "Debug"

	t, err := CompileTemplate("debug.html")
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