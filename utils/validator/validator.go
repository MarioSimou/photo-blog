package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Middleware struct{}

func (m Middleware) ValidateCreateUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		bf, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(bf))
		var body interface{}
		json.NewDecoder(r.Body).Decode(&body)
		data, _ := body.(map[string]string)

		fmt.Println(data)
		if _, ok := data["username"]; !ok {
			http.Error(w, "Unable to parse user data", 400)
			return
		}

		if _, ok := data["email"]; !ok {
			http.Error(w, "Unable to parse user data", 400)
			return
		}

		if _, ok := data["password"]; !ok {
			http.Error(w, "Unable to parse user data", 400)
			return
		}

		fmt.Println("validate user input")

		next(w, r, p)
	}
}

func (m Middleware) ValidateRequest(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//HTTP/x.x 406 Not Acceptable
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

		bf, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(bf))
		next(w, r, p)
	}
}
