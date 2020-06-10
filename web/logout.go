package web

import (
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us, err := store.Get(r, "session-key")
	if err != nil {
		logger.Warningf("Could not open session: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set user to nil
	us.Values["user"] = nil
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}