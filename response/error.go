package response

import "net/http"

func Error(code int, msg string) Response {
	return new(code, err, "error", msg)
}

func InternalServerError() Response {
	msg := "the server encontered a problem and could not process your request"
	return Error(http.StatusInternalServerError, msg)
}
