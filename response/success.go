package response

import (
	"net/http"
)

func Success(code int, key string, data any) Response {
	return new(code, success, key, data)
}

func OK(key string, data any) Response {
	return Success(http.StatusOK, key, data)
}

func Created(key string, data any) Response {
	return Success(http.StatusCreated, key, data)
}
