package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"projects/users-auth-api/controllers"
	"projects/users-auth-api/utils"
	"projects/users-auth-api/utils/httpcodes"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var c *controllers.Controller
var m Middleware

func customRoute(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintln(w, "Ping")
}

func checkHeader(w *httptest.ResponseRecorder, k string, e string, t *testing.T) {
	if v := w.Header().Get(k); v != e {
		t.Errorf("Should have returned a header value of %v rather than %v", e, v)
	}
}
func checkStatusCode(res *http.Response, e int, t *testing.T) {
	if res.StatusCode != e {
		t.Errorf("Should have returned a status code of %d rather than %d", e, res.StatusCode)
	}
}
func parseResponseBody(res *http.Response) *httpcodes.Representation {
	var repr httpcodes.Representation
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &repr)
	return &repr
}

func init() {
	u := utils.Utils{}
	u.LoadDotEnv("../../.test.env")
	m = Middleware{Utils: &u}
	mcli := u.ConnectDatabase(os.Getenv("MONGO_URI"), os.Getenv("DB_NAME"))
	c = controllers.NewController(mcli, &u)
}

func TestCreateUser(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"username": "paul","email":"paul@gmail.com","password":"12345678"}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateCreateUser(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 200, t)
	checkHeader(w, "Content-Type", "application/json", t)
}

func TestCreateUserInvalidUsername(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"username": "","email":"paul@gmail.com","password":"12345678"}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateCreateUser(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 400, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 400 {
		t.Errorf("Should have returned a status code of 400 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "Invalid Username" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid Username", repr.Message)
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"username": "paul","email":"","password":"12345678"}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateCreateUser(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 400, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 400 {
		t.Errorf("Should have returned a status code of 400 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "Invalid Email" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid Email", repr.Message)
	}
}

func TestCreateUserInvalidPassword(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"username": "paul","email":"paul@gmail.com","password":""}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateCreateUser(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 400, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 400 {
		t.Errorf("Should have returned a status code of 400 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "Invalid Password" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid Password", repr.Message)
	}
}

func TestCreateUserEmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	m.ValidateCreateUser(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 400, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 400 {
		t.Errorf("Should have returned a status code of 400 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
}

func TestValidateSignIn(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"email":"paul@gmail.com", "password":"12345678"}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateSignIn(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 200, t)
	checkHeader(w, "Content-Type", "application/json", t)
}

func TestValidateSignInInvalidEmail(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"email":"", "password":"12345678"}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateSignIn(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 401, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 401 {
		t.Errorf("Should have returned a status code of 401 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "Invalid Request Body" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid Request Body", repr.Message)
	}
}

func TestValidateSignInInvalidPassword(t *testing.T) {
	w := httptest.NewRecorder()
	body := []byte(`{"email":"paul@gmail.com", "password":""}`)
	r := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	m.ValidateSignIn(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 401, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 401 {
		t.Errorf("Should have returned a status code of 401 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "Invalid Request Body" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid Request Body", repr.Message)
	}
}

func TestValidateHeaders(t *testing.T) {

}
