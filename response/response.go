package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Header http.Header
	Code   int
	Body   any
}

var checkCode = func(c int) {
	if http.StatusText(c) == "" {
		panic(fmt.Sprintf("response code %d is unknown", c))
	}
}

func New(code int) *Response {
	checkCode(code)
	return &Response{
		Header: make(http.Header, 0),
		Code:   code,
	}
}

func (rw *Response) JSON(w http.ResponseWriter, body any) *Response {
	f := json.MarshalIndent
	rw.Body = body
	if err := write(w, f, rw); err != nil {
		// As developer we need to be sure our models are serializable
		// if not it will panic and return 500 to the client
		// when the server middleware recover panic raises
		panic(err)
	}
	return rw
}
