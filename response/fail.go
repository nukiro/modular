package response

import (
	"fmt"
	"net/http"
)

// 4xx responses
func Fail(code int, errors any) Response {
	return new(code, fail, "error", errors)
}

func BadRequest(errors any) Response {
	return Fail(http.StatusBadRequest, errors)
}

func NotFound() Response {
	msg := "the requested resource could not be found"
	return Fail(http.StatusNotFound, msg)
}

func MethodNotAllowed(method string) Response {
	msg := fmt.Sprintf("the %s method is not supported for this resource", method)
	return Fail(http.StatusMethodNotAllowed, msg)
}

func Conflict() Response {
	msg := "unable to update the record due to an edit conflict, please try again"
	return Fail(http.StatusConflict, msg)
}

func UnprocessableEntity(errors any) Response {
	return Fail(http.StatusUnprocessableEntity, errors)
}
