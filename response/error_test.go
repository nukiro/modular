package response

import "testing"

func TestInternalServerErrorResponse(t *testing.T) {
	rw := InternalServerError()
	assertResponseStatusCode(t, rw, "500 Internal Server Error")
}
