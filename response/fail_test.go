package response

import "testing"

func TestNotFoundResponse(t *testing.T) {
	rw := NotFound()
	assertResponseStatusCode(t, rw, "404 Not Found")
}

func TestBadRequestResponse(t *testing.T) {
	rw := BadRequest("errors")
	assertResponseStatusCode(t, rw, "400 Bad Request")
}

func TestMethodNotAllowedResponse(t *testing.T) {
	rw := MethodNotAllowed("GET")
	assertResponseStatusCode(t, rw, "405 Method Not Allowed")
}

func TestConflictResponse(t *testing.T) {
	rw := Conflict()
	assertResponseStatusCode(t, rw, "409 Conflict")
}

func TestUnprocessableEntityResponse(t *testing.T) {
	rw := UnprocessableEntity("errors")
	assertResponseStatusCode(t, rw, "422 Unprocessable Entity")
}
