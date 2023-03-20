package services

import "net/http"

func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/ident", app.BasicAuth(app.logQuery))
	mux.HandleFunc("/addIncident", app.BasicAuth(app.addIncident))
	mux.HandleFunc("/showIncident", app.BasicAuth(app.showIncident))
	return mux
}
