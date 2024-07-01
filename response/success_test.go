package response

import "testing"

func TestOkResponse(t *testing.T) {
	rw := OK("key", "data")
	assertResponseStatusCode(t, rw, "200 OK")
}

func TestCreatedResponse(t *testing.T) {
	rw := Created("key", "data")
	assertResponseStatusCode(t, rw, "201 Created")
}
