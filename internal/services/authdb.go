package services

import (
	"golang.org/x/crypto/bcrypt"
)

// verifyUserPass verifies that username/password is a valid pair matching
// our userPasswords "database".
func verifyUserPass(app *application, username, password string) bool {
	user, err := app.Rep.GetUserByNameRep(username)
	if err != nil {

		return false
	}
	wantPass := []byte(user.Pass)

	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}

/*func UserRule(app *application, username string) string {
	user, err := app.modelsDB.GetUserByName(username)
	if err != nil {

		return ""
	}
	rule, err := app.modelsDB.GetUserRules(user.Id)
	if err != nil {

		return ""
	}
	return rule.Name
}*/
