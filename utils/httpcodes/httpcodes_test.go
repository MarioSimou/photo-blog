package httpcodes

import "testing"

var repr Representation

func init() {
	repr = Representation{Message: ""}
}

func checkStatusCode(r *Representation, e int, t *testing.T) {
	if r.Status != e {
		t.Errorf("Should have returned a status code of %v rather than %v", e, r.Status)
	}
}
func checkSuccess(r *Representation, e bool, t *testing.T) {
	if r.Success != e {
		t.Errorf("Should have returned a boolean flag of '%v' rather than %v", e, r.Status)
	}
}

func TestBadRequest(t *testing.T) {
	r := repr.BadRequest()
	checkStatusCode(&r, 400, t)
	checkSuccess(&r, false, t)
}

func TestUnauthorized(t *testing.T) {
	r := repr.Unauthorized()
	checkStatusCode(&r, 401, t)
	checkSuccess(&r, false, t)
}

func TestForbidden(t *testing.T) {
	r := repr.Forbidden()
	checkStatusCode(&r, 403, t)
	checkSuccess(&r, false, t)
}
func TestNotFound(t *testing.T) {
	r := repr.NotFound()
	checkStatusCode(&r, 404, t)
	checkSuccess(&r, false, t)
}
func TestUnsupportedMediaType(t *testing.T) {
	r := repr.UnsupportedMediaType()
	checkStatusCode(&r, 415, t)
	checkSuccess(&r, false, t)
}

func TestInternalServerError(t *testing.T) {
	r := repr.InternalServerError()
	checkStatusCode(&r, 500, t)
	checkSuccess(&r, false, t)
}
func TestOk(t *testing.T) {
	r := repr.Ok()
	checkStatusCode(&r, 200, t)
	checkSuccess(&r, true, t)
}
func TestCreated(t *testing.T) {
	r := repr.Created()
	checkStatusCode(&r, 201, t)
	checkSuccess(&r, true, t)
}
