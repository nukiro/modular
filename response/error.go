package response

import "net/http"

var internalServerErrorMsg = "the server encontered a problem and could not process your request"

func Error(code int, msg string) Response {
	return new(code, err, "error", msg)
}

func InternalServerError() Response {
	return Error(http.StatusInternalServerError, internalServerErrorMsg)
}
