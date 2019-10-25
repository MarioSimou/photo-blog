package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/MarioSimou/photo-blog-in-golang/models"
	"github.com/MarioSimou/photo-blog-in-golang/utils"
	"github.com/MarioSimou/photo-blog-in-golang/utils/httpcodes"
	"github.com/julienschmidt/httprouter"
)

type MiddlewareHandler func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{})

func Handler(next MiddlewareHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, r, p)
	}
}

type Middleware struct {
	Utils *utils.Utils
}

func (m Middleware) ValidateCreateUser(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		var body models.User
		json.NewDecoder(r.Body).Decode(&body)

		if t := body.ValidateUsername(); !t {
			// HTTP/x.x 400 Bad Request
			json.NewEncoder(w).Encode(httpcodes.Response{Message: "Invalid Username"}.BadRequest())
			return
		}
		if t := body.ValidateEmail(); !t {
			// HTTP/x.x 400 Bad Request
			json.NewEncoder(w).Encode(httpcodes.Response{Message: "Invalid Email"}.BadRequest())
			return
		}
		if t := body.ValidatePassword(); !t {
			// HTTP/x.x 400 Bad Request
			json.NewEncoder(w).Encode(httpcodes.Response{Message: "Invalid Password"}.BadRequest())
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

		// HTTP/x.x 400 Bad Request
		json.NewEncoder(w).Encode(httpcodes.Response{Message: "Invalid Request Body"}.BadRequest())
	}
}

func (m Middleware) ValidateRequest(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		// HTTP/x.x 406 Not Acceptable
		if a := r.Header.Get("Accept"); a != "*/*" && a != "application/json" {
			json.NewEncoder(w).Encode(httpcodes.Response{Message: "Only JSON representations are supported"}.NotAcceptable())
			return
		}
		// HTTP/x.x 415 Unsupported Media Type
		if m := r.Method; m == http.MethodPost || m == http.MethodPut {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				json.NewEncoder(w).Encode(httpcodes.Response{Message: "A MIME type of application/json is only accepted"}.UnsupportedMediaType())
				return
			}
		}

		next(w, r, p)
	}
}

func (m Middleware) Authorization(next MiddlewareHandler) MiddlewareHandler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
		auth := r.Header.Get("Authorization")
		t := strings.Replace(auth, "Bearer ", "", 1)
		if t == "" || auth == "" {
			// HTTP/x.x 401 Unauthorized
			json.NewEncoder(w).Encode(httpcodes.Response{Message: "Invalid user token"}.Unauthorized())
			return
		}

		if _, ok := m.Utils.VerifyToken([]byte(t), os.Getenv("JWT_SECRET")); ok {
			// custom header that includes the JWT token
			payload := m.Utils.ExtractPayload(t)
			next(w, r, p, payload)
		} else {
			// HTTP/x.x 401 Unauthorized
			json.NewEncoder(w).Encode(httpcodes.Response{Message: "Invalid user token"}.Unauthorized())
			return
		}
	}
}
