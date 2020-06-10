package web

import (
	"encoding/json"
	"net/http"
	"phsite/models"
)

func HandleOauthCallback(w http.ResponseWriter, r *http.Request) {
	// Init Session
	us, err := store.Get(r, "session-key")
	if err != nil {
		logger.Warningf("Could not open session: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get state
	val := us.Values["oauth-state"]
	var state string
	var ok bool
	if state, ok = val.(string); !ok {
		// redirect home page if no login-redirect
		logger.Warningf("Invalid State")
		http.Error(w, "state invalid", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state {
		logger.Warningf("States don't match")
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// display error
	if r.URL.Query().Get("error") != "" {
		http.Error(w, r.URL.Query().Get("error_description"), http.StatusForbidden)
		return
	}

	// process returned oauth data
	user, err := processCallback(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not logged in", http.StatusForbidden)
		return
	}

	// Insert into database or update existing record0
	err = user.Upsert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save user to session
	us.Values["oauth-state"] = nil
	us.Values["user"] = user
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to last page
	val = us.Values["login-redirect"]
	var loginRedirect string
	if loginRedirect, ok = val.(string); !ok {
		// redirect home page if no login-redirect
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, loginRedirect, http.StatusFound)
	return
}

func processCallback(r *http.Request) (*models.User, error) {
	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		return nil, err
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, err
	}
	idToken, err := oauth2Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	IDTokenClaims := new(json.RawMessage)
	if err := idToken.Claims(&IDTokenClaims); err != nil {
		return nil, err
	}

	logger.Debugf("Response From OAUTH: %s", IDTokenClaims)

	user := models.User{}
	if err := json.Unmarshal(*IDTokenClaims, &user); err != nil {
		return nil, err
	}

	return &user, nil
}