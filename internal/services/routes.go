package services

import "net/http"

func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.PanicRecovery(app.home))
	mux.HandleFunc("/ident", app.PanicRecovery(app.BasicAuth(app.logQuery)))
	mux.HandleFunc("/addIncident", app.PanicRecovery(app.TokenAuth(app.addIncident)))
	mux.HandleFunc("/showIncident", app.PanicRecovery(app.TokenAuth(app.showIncident)))

	//handler := app.PanicRecovery(mux.HandleFunc())
	return mux
}
