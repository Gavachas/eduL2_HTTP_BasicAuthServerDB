package services

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// UserContextKey is the key in a request's context used to check if the request
// has an authenticated user. The middleware will set the value of this key to
// the username, if the user was properly authenticated with a password.
const UserContextKey = "user"

// BasicAuth is middleware that verifies the request has appropriate basic auth
// set up with a user:password pair verified by authdb.
func (app *application) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//if verifyUserToken(app, w, req) {
		//	next.ServeHTTP(w, req)
		//} else {
		user, pass, ok := req.BasicAuth()

		if ok && verifyUserPass(app, user, pass) {
			//fmt.Fprintf(w, "Hello user!\n")
			usertoken := GenerateToken(user)
			//err := app.Session.Set(user, usertoken)
			err := app.Session.Set(usertoken, user)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				w.Header().Set("Token", usertoken)
				livingTime := 60 * time.Minute
				expiration := time.Now().Add(livingTime)
				//кука будет жить 1 час
				cookie := http.Cookie{Name: "token", Value: url.QueryEscape(usertoken), Expires: expiration}
				http.SetCookie(w, &cookie)
				newctx := context.WithValue(req.Context(), UserContextKey, user)
				next.ServeHTTP(w, req.WithContext(newctx))
			}
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
	//}
}
func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}
func (app *application) TokenAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		tokenCoookie, err := readCookie("token", req)
		if err != nil {
			app.logger.Error("не найден токен")
			http.Redirect(w, req, "/ident", http.StatusSeeOther)
			return
		}
		if tokenCoookie != "" {
			userTokenSession, err := app.Session.Get(tokenCoookie)
			if err != nil {
				app.logger.Error("не найден токен в редис")
				http.Redirect(w, req, "/ident", http.StatusSeeOther)
				return
			}
			newctx := context.WithValue(req.Context(), UserContextKey, userTokenSession)
			next.ServeHTTP(w, req.WithContext(newctx))
			return
		}

		userContext, ok := req.Context().Value("user").(string)
		if !ok {
			app.logger.Error("не найден user")
			http.Redirect(w, req, "/ident", http.StatusSeeOther)
			return
		}
		userToken := w.Header().Get("Token")
		userTokenSession, err := app.Session.Get(userContext)
		if err != nil {
			app.logger.Error("не найден токен")
			http.Redirect(w, req, "/ident", http.StatusSeeOther)
			return
		}
		tokenSessionS := userTokenSession.(string)
		if tokenSessionS == userToken {
			//fmt.Fprintf(w, "Hello user!\n")

			next.ServeHTTP(w, req)
		} else {
			//w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			//http.Error(w, "Unauthorized", http.StatusUnauthorized)
			http.Redirect(w, req, "/ident", http.StatusSeeOther)
		}
	}
}
func (app *application) PanicRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				if errString, ok := err.(string); ok {
					app.serverError(w, errors.New(errString))
				} else {
					app.serverError(w, errors.New("Panic recovery error!"))
				}

			}
		}()
		next.ServeHTTP(w, req)
	}
}
