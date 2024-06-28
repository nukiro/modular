package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Response interface {
	Header(key string, value string)
	Write(w http.ResponseWriter) int
}

type payload map[string]any
type response struct {
	headers http.Header
	code    int
	payload
}

type result string

var (
	success result = "success"
	fail    result = "fail"
	err     result = "error"
)

func new(code int, result result, key string, data any) *response {
	if http.StatusText(code) == "" {
		panic(fmt.Sprintf("response code %d is unknown", code))
	}

	if key == "" {
		panic("response key cannot be empty")
	}

	return &response{
		headers: nil,
		code:    code,
		payload: payload{
			"time":   time.Now().Unix(),
			"status": strings.ToLower(http.StatusText(code)),
			"result": result,
			key:      data,
		},
	}
}

func (r *response) writeError(w http.ResponseWriter) {
	message := "the server encontered a problem and could not process your request"
	r.code = http.StatusInternalServerError
	errorResponse := new(r.code, "error", "error", message)
	if err := errorResponse.write(w); err != nil {
		w.WriteHeader(500)
		return
	}
}

func (r *response) Header(key string, value string) {
	if r.headers == nil {
		r.headers = make(http.Header)
	}
	r.headers.Set(string(key), value)
}

func (r *response) Write(w http.ResponseWriter) int {
	if err := r.write(w); err != nil {
		r.writeError(w)
	}
	return r.code
}

// write sends a JSON response to the client.
func (r *response) write(w http.ResponseWriter) error {
	// MarshalIndent adds whitespaces to the encoded JSON.
	// No line prefix ("") and two white spaces indents ("\t") for each element.
	js, err := json.MarshalIndent(r.payload, "", "  ")
	if err != nil {
		return err
	}

	// Append a new line making it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point it is safe to add any headers as we know that we will not
	// encounter any more errors before writing the response.
	// Custom headers value pass by param method.
	for key, value := range r.headers {
		w.Header()[key] = value
	}
	// Response Content Type
	w.Header().Set("Content-Type", "application/json")
	// Response Status Code
	w.WriteHeader(r.code)

	w.Write(js)
	return nil
}

type serializer func(v any, prefix, indent string) ([]byte, error)

func serialize(r *response, f serializer) ([]byte, error) {
	// MarshalIndent adds whitespaces to the encoded JSON.
	// No line prefix ("") and two white spaces indents ("\t") for each element.
	js, err := f(r.payload, "", "  ")
	if err != nil {
		return nil, err
	}

	// Append a new line making it easier to view in terminal applications.
	js = append(js, '\n')
	return js, nil
}

// Write sends a JSON response to the client.
func write(r *response, w http.ResponseWriter, f serializer) error {
	js, err := serialize(r, f)
	if err != nil {
		return err
	}

	// At this point it is safe to add any headers as we know that we will not
	// encounter any more errors before writing the response.
	// Custom headers value pass by param method.
	for key, value := range r.headers {
		w.Header()[key] = value
	}
	// Response Content Type
	w.Header().Set("Content-Type", "application/json")
	// Response Status Code
	w.WriteHeader(r.code)

	w.Write(js)
	return nil
}
