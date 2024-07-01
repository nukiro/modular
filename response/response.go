package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Response interface {
	Status() string
	Body() any
	Header() http.Header
	Write(w http.ResponseWriter) Response
}

type body map[string]any
type response struct {
	header     http.Header
	statusCode int
	body
}

type result string

var (
	success result = "success"
	fail    result = "fail"
	err     result = "error"
)

func checkStatusCode(c int) {
	if http.StatusText(c) == "" {
		panic(fmt.Sprintf("response code %d is unknown", c))
	}
}

func new(statusCode int, result result, key string, data any) *response {
	checkStatusCode(statusCode)

	if key == "" {
		panic("response key param cannot be empty")
	}

	if data == "" {
		panic("response data param cannot be empty")
	}

	return &response{
		header:     make(http.Header, 0),
		statusCode: statusCode,
		body: body{
			"time":   time.Now().Unix(),
			"status": strings.ToLower(http.StatusText(statusCode)),
			"result": result,
			key:      data,
		},
	}
}

func (rw *response) Status() string {
	return fmt.Sprintf("%d %s", rw.statusCode, http.StatusText(rw.statusCode))
}

func (rw *response) Body() any {
	return rw.body
}

func (rw *response) Header() http.Header {
	return rw.header
}

func (rw *response) Write(w http.ResponseWriter) Response {
	f := json.MarshalIndent
	if err := write(w, f, rw); err != nil {
		return writeError(w, f)
	}
	return rw
}
