package web

import (
	"net/http"
	"phsite/util"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us, err := store.Get(r, "session-key")
	if err != nil {
		logger.Warningf("Could not open session: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newState := util.RandString(16)
	us.Values["oauth-state"] = newState
	err = us.Save(r, w)
	if err != nil {
		logger.Warningf("Could not save session: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	http.Redirect(w, r, oauth2Config.AuthCodeURL(newState), http.StatusFound)
}