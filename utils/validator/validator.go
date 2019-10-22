package validator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"../../models"
	"github.com/julienschmidt/httprouter"
)

type Middleware struct{}

func (m Middleware) ValidateCreateUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var body models.User
		bf, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bf))

		e := json.Unmarshal(bf, &body)
		if e != nil {
			http.Error(w, "Unable to parse user data", 400)
			return
		}
		if body.Username == "" {
			http.Error(w, "Unable to parse user data", 400)
			return
		}

		if body.Email == "" {
			http.Error(w, "Unable to parse user data", 400)
			return
		}

		if body.Password == "" {
			http.Error(w, "Unable to parse user data", 400)
			return
		}

		if len(body.Password) < 8 {
			http.Error(w, "User password should be more than 8 characters", 400)
			return
		}

		next(w, r, p)
	}
}

func (m Middleware) ValidateRequest(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// HTTP/x.x 406 Not Acceptable
		if a := r.Header.Get("Accept"); a != "*/*" && a != "application/json" {
			http.Error(w, "Only JSON representation are supported", 406)
			return
		}

		// HTTP/x.x 415 Unsupported Media Type
		if m := r.Method; m == http.MethodPost || m == http.MethodPut {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				http.Error(w, "A MIME type of application/json is only accepted", 415)
				return
			}
		}

		next(w, r, p)
	}
}
