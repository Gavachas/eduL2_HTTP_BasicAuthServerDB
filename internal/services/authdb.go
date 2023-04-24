package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(user), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash :", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

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
func verifyUserToken(app *application, w http.ResponseWriter, req *http.Request) bool {

	userContext, ok := req.Context().Value("user").(string)
	if !ok {
		//app.serverError(w, errors.New("не найден user"))
		return false
	}
	userToken := w.Header().Get("Token")
	userTokenSession, err := app.Session.Get(userContext)
	if err != nil || userTokenSession == nil {
		//app.serverError(w, errors.New("не найден токен"))
		return false
	}
	tokenSessionS := userTokenSession.(string)

	return tokenSessionS == userToken
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
