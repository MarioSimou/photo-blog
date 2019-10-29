// Package middlewares provides functionalities that validate the request before its actual execution
package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/MarioSimou/authAPI/internal/models"
	"github.com/MarioSimou/authAPI/internal/utils"
	"github.com/MarioSimou/authAPI/internal/utils/httpcodes"

	"github.com/julienschmidt/httprouter"
)

// MiddlewareHandler is a custom type that extends the capabilities of httprouter.Handle
type MiddlewareHandler func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{})

// Handler is a closure for ht
func Handler(next MiddlewareHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, r, p)
	}
}

// Middleware is a custom type that accepts the Utilities type from Utilities package
type Middleware struct {
	Utils *utils.Utils
}

// ValidateCreateUser validates the request body when a user is created
func (m Middleware) ValidateCreateUser(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		var body models.User
		json.NewDecoder(r.Body).Decode(&body)

		if t := body.ValidateUsername(); !t {
			// HTTP/x.x 400 Bad Request
			httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid Username"}.BadRequest())
			return
		}
		if t := body.ValidateEmail(); !t {
			// HTTP/x.x 400 Bad Request
			httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid Email"}.BadRequest())
			return
		}
		if t := body.ValidatePassword(); !t {
			// HTTP/x.x 400 Bad Request
			httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid Password"}.BadRequest())
			return
		}
		body.Password = m.Utils.HashPassword(body.Password)
		body.ValidateRole()

		// updates the content of the request body
		nB, _ := json.Marshal(body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(nB))
		r.Body.Close()
		next(w, r, p)
	}
}

// ValidateSignIn checks the credentials of a user. If the calidation failsm it returns an HTTP 401 Unauthorized.
func (m Middleware) ValidateSignIn(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		var body models.LoginUser
		json.NewDecoder(r.Body).Decode(&body)

		if body.ValidateEmail() && body.ValidatePassword() {
			j, _ := json.Marshal(body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(j))
			r.Body.Close()
			next(w, r, p)
			return
		}

		// HTTP/x.x 401 Unauthorized
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid Request Body"}.Unauthorized())
		return
	}
}

// ValidateRequest checks the Request Headers(Accept and Content-Type) of a request, which shows that the API
// either accept or returns data in a JSON format
func (m Middleware) ValidateRequest(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		// HTTP/x.x 406 Not Acceptable
		if a := r.Header.Get("Accept"); a != "*/*" && a != "application/json" {
			httpcodes.ResponseError(w, httpcodes.Representation{Message: "Only JSON representations are supported"}.NotAcceptable())
			return
		}
		// HTTP/x.x 415 Unsupported Media Type
		if m := r.Method; m == http.MethodPost || m == http.MethodPut {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				httpcodes.ResponseError(w, httpcodes.Representation{Message: "A MIME type of application/json is only accepted"}.UnsupportedMediaType())
				return
			}
		}

		next(w, r, p)
	}
}

// Authorization checks the credentials of a user. A user needs to use a valid JWT token.
func (m Middleware) Authorization(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		auth := r.Header.Get("Authorization")
		t := strings.Replace(auth, "Bearer ", "", 1)
		if t == "" || auth == "" {
			// HTTP/x.x 401 Unauthorized
			httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid user token"}.Unauthorized())
			return
		}

		if _, ok := m.Utils.VerifyToken([]byte(t), os.Getenv("JWT_SECRET")); ok {
			// custom header that includes the JWT token
			payload := m.Utils.ExtractPayload(t)
			next(w, r, p, payload)
		} else {
			// HTTP/x.x 401 Unauthorized
			httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid user token"}.Unauthorized())
			return
		}
	}
}
