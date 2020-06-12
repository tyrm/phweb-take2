package web

import "net/http"

type HomeTemplate struct {
	templateCommon
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmplVars :=  HomeTemplate{}
	_, err := initSession(w, r, &tmplVars)

	tmplVars.PageTitle = "Home"

	t, err := CompileTemplate("home.html")
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
