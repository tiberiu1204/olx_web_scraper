package middleware

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tiberiu1204/olx_web_scraper/api"
	"github.com/tiberiu1204/olx_web_scraper/internal/tools"
)

var InvalidUsername = errors.New("Invalid username.")
var InvalidToken = errors.New("Invalid token.")
var UnAuthorizedError = errors.New("Invalid username and password combination.")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string = r.URL.Query().Get("username")
		var token string = r.Header.Get("Authorization")

		if username == "" {
			log.Error(InvalidUsername)
			api.RequestErrorHandler(w, InvalidUsername)
			return
		}

		if token == "" {
			log.Error(InvalidToken)
			api.RequestErrorHandler(w, InvalidToken)
			return
		}

		var database *tools.DatabaseInterface
		database, err := tools.NewDatabase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		var loginDetails *tools.LoginDetails = (*database).GetUserLoginDetails(username)

		if loginDetails == nil || (token != loginDetails.AuthToken) {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
