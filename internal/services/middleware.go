package services

import (
	"context"
	"fmt"
	"net/http"
)

// UserContextKey is the key in a request's context used to check if the request
// has an authenticated user. The middleware will set the value of this key to
// the username, if the user was properly authenticated with a password.
const UserContextKey = "user"

// BasicAuth is middleware that verifies the request has appropriate basic auth
// set up with a user:password pair verified by authdb.
func (app *application) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()

		if ok && verifyUserPass(app, user, pass) {
			fmt.Fprintf(w, "Hello user!\n")
			newctx := context.WithValue(req.Context(), UserContextKey, user)
			next.ServeHTTP(w, req.WithContext(newctx))
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
