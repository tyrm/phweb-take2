package web

import "net/http"

type MastodonTemplate struct {
	templateCommon
}

func HandleMastodon(w http.ResponseWriter, r *http.Request) {
	tmplVars :=  HomeTemplate{}
	_, err := initSession(w, r, &tmplVars)

	tmplVars.PageTitle = "Mastodon"

	t, err := CompileTemplate("mastodon.html")
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
