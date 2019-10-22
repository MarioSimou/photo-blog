package validator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	utils ".."
	"../../models"
	"github.com/julienschmidt/httprouter"
)

type Middleware struct{}

func (m Middleware) ValidateCreateUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var body models.User
		// the request body is updated since it can only read once
		bf, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bf))

		e := json.Unmarshal(bf, &body)
		if e != nil {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Unable to parse data",
			})
			return
		}
		if body.Username == "" {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Invalid username",
			})
			return
		}

		if body.Email == "" {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Invalid user email",
			})
			return
		}

		if body.Password == "" {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "Invalid user password",
			})
			return
		}

		if len(body.Password) < 8 {
			json.NewEncoder(w).Encode(utils.Response{
				Status:  400,
				Success: false,
				Message: "The user password should be more than 8 character",
			})
			return
		}

		next(w, r, p)
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
			}
		}

		next(w, r, p)
	}
}
