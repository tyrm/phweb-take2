package web

import (
	"net/http"
	"phsite/models"
	"time"
)

func MiddlewareRequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Init Session
		us, err := store.Get(r, "session-key")
		if err != nil {
			logger.Infof("got %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve our struct and type-assert it
		val := us.Values["user"]
		var user = models.User{}
		var ok bool
		if user, ok = val.(models.User); !ok {
			// Save current page
			us.Values["login-redirect"] = r.URL.Path
			err = us.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		now := time.Now()
		if user.ExpiresAt < now.Unix() {
			// Save current page
			us.Values["login-redirect"] = r.URL.Path
			err = us.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}