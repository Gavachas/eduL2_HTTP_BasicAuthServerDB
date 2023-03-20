package services

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Hello guest!\n")
}
func (app *application) logQuery(w http.ResponseWriter, req *http.Request) {
	userContext := req.Context().Value("user")
	user, err := app.Rep.GetUserByNameRep(userContext.(string))
	if err != nil {
		app.serverError(w, err)
		return
	}
	userRule, err := app.Rep.GetUserRulesRep(user.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if userRule.Name == "admin" {
		fmt.Fprintf(w, "Hello admin! You can del incidents\n ")
	} else if userRule.Name == "user" {
		fmt.Fprintf(w, "Hello admin! You can add incidents\n ")

	}

}
func (app *application) addIncident(w http.ResponseWriter, req *http.Request) {
	userContext := req.Context().Value("user")
	user, err := app.Rep.GetUserByNameRep(userContext.(string))
	if err != nil {
		app.serverError(w, err)
		return
	}
	userRule, err := app.Rep.GetUserRulesRep(user.Id)
	if err != nil {

		app.serverError(w, err)
		return
	}
	if userRule.Name == "admin" {

		inc, err := app.Rep.InsertIncidetRep("Admin make inc", user.Id)

		if err != nil {
			app.serverError(w, err)
			return
		}

		fmt.Fprintf(w, "Incident ID : %v \n", inc)
	} else if userRule.Name == "user" {

		inc, err := app.Rep.InsertIncidetRep("User make inc", user.Id)

		if err != nil {
			app.serverError(w, err)
			return
		}
		fmt.Fprintf(w, "Hello user! You can add incidents\n ")
		fmt.Fprintf(w, "Incident ID : %v \n", inc)
	}

}
func (app *application) showIncident(w http.ResponseWriter, req *http.Request) {
	userContext := req.Context().Value("user")
	id, err := strconv.Atoi(req.URL.Query().Get("id"))

	if err != nil || id < 0 {
		http.NotFound(w, req)
		return
	}

	user, err := app.Rep.GetUserByNameRep(userContext.(string))
	if err != nil {
		app.serverError(w, err)
		return
	}
	userRule, err := app.Rep.GetUserRulesRep(user.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if userRule.Name == "admin" {
		inc, err := app.Rep.GetIncidentRep(id)

		if err != nil {
			app.serverError(w, err)
			return
		}
		fmt.Fprintf(w, "Incident  : %v \n", inc)
	} else if userRule.Name == "user" {
		fmt.Fprintf(w, "Access denied \n")
	}

}
