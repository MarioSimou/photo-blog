package validator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/MarioSimou/photo-blog-in-golang/models"
	"github.com/MarioSimou/photo-blog-in-golang/utils"
	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	Utils *utils.Utils
}

func (m Middleware) ValidateCreateUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var body models.User
		json.NewDecoder(r.Body).Decode(&body)

		if t := body.ValidateUsername(); !t {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Invalid username",
			})
			return
		}
		if t := body.ValidateEmail(); !t {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Invalid user email",
			})
			return
		}
		if t := body.ValidatePassword(); !t {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Invalid user password",
			})
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

func (m Middleware) ValidateLoginUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var body models.LoginUser
		json.NewDecoder(r.Body).Decode(&body)
		if body.ValidateEmail() && body.ValidatePassword() {
			j, _ := json.Marshal(body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(j))
			r.Body.Close()
			next(w, r, p)
			return
		}

		json.NewEncoder(w).Encode(utils.Response{
			Status:  400,
			Success: false,
			Message: "Invalid request body",
		})
	}
}

func (m Middleware) ValidateRequest(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// HTTP/x.x 406 Not Acceptable
		if a := r.Header.Get("Accept"); a != "*/*" && a != "application/json" {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  406,
				Success: false,
				Message: "Only JSON representation are supported",
			})
			return
		}
		// HTTP/x.x 415 Unsupported Media Type
		if m := r.Method; m == http.MethodPost || m == http.MethodPut {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				json.NewEncoder(w).Encode(utils.Response{
					Status:  415,
					Success: false,
					Message: "A MIME type of application/json is only accepte",
				})
				return
			}
		}

		next(w, r, p)
	}
}

func (m Middleware) Authorization(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		auth := r.Header.Get("Authorization")
		t := strings.Replace(auth, "Bearer ", "", 1)
		if t == "" || auth == "" {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  401,
				Success: false,
				Message: "Invalid user token",
			})
			return
		}

		if _, ok := m.Utils.VerifyToken([]byte(t), os.Getenv("JWT_SECRET")); ok {
			next(w, r, p)
		} else {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  401,
				Success: false,
				Message: "Invalid user token",
			})
		}
	}
}
