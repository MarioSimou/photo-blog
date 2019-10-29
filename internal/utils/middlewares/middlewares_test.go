package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/MarioSimou/authAPI/internal/controllers"
	"github.com/MarioSimou/authAPI/internal/models"
	"github.com/MarioSimou/authAPI/internal/utils"
	"github.com/MarioSimou/authAPI/internal/utils/httpcodes"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var c *controllers.Controller
var m Middleware
var u utils.Utils

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
	u = utils.Utils{}
	u.LoadDotEnv("../../../configs/.test.env")
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

func TestValidateRequestAcceptHeaderXML(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	r.Header.Set("Accept", "application/xml")
	m.ValidateRequest(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 406, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 406 {
		t.Errorf("Should have returned a status code of 406 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "Only JSON representations are supported" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Only JSON representations are supported", repr.Message)
	}
}

func TestValidateRequestContentType(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	r.Header.Set("Accept", "application/json")
	m.ValidateRequest(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 415, t)
	checkHeader(w, "Content-Type", "application/json", t)
	repr := parseResponseBody(res)
	if repr.Status != 415 {
		t.Errorf("Should have returned a status code of 415 rather than %v", repr.Status)
	}
	if repr.Success {
		t.Errorf("Should have returned a 'false' flag")
	}
	if repr.Message != "A MIME type of application/json is only accepted" {
		t.Errorf("Should have returned an error message of %v rather than %v", "A MIME type of application/json is only accepted", repr.Message)
	}
}

func TestValidateRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	m.ValidateRequest(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 200, t)
	checkHeader(w, "Content-Type", "application/json", t)
}

func TestAuthorizationWithoutToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	m.Authorization(customRoute)(w, r, nil)

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
	if repr.Message != "Invalid user token" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid user token", repr.Message)
	}
}

func TestAuthorizationInvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	user := models.User{Id: &userId, Email: "paul@gmail.com"}
	// invalid token
	token, _ := u.GenerateToken(user, "randomsecret", time.Hour)
	r.Header.Set("Authorization", "Bearer "+string(token))
	m.Authorization(customRoute)(w, r, nil)

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
	if repr.Message != "Invalid user token" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid user token", repr.Message)
	}
}

func TestAuthorizationExpiredToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	user := models.User{Id: &userId, Email: "paul@gmail.com"}
	// invalid token
	token, _ := u.GenerateToken(user, os.Getenv("JWT_SECRET"), 0*time.Second)
	time.Sleep(1 * time.Second)
	r.Header.Set("Authorization", "Bearer "+string(token))
	m.Authorization(customRoute)(w, r, nil)

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
	if repr.Message != "Invalid user token" {
		t.Errorf("Should have returned an error message of %v rather than %v", "Invalid user token", repr.Message)
	}
}

func TestAuthorizationValidToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users", nil)
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	user := models.User{Id: &userId, Email: "paul@gmail.com"}
	// invalid token
	token, _ := u.GenerateToken(user, os.Getenv("JWT_SECRET"), time.Hour)
	r.Header.Set("Authorization", "Bearer "+string(token))
	m.Authorization(customRoute)(w, r, nil)

	res := w.Result()
	checkStatusCode(res, 200, t)
	checkHeader(w, "Content-Type", "application/json", t)
}
