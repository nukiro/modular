package response

import (
	"net/http"
)

type Serializer func(v any, prefix, indent string) ([]byte, error)

var serialize = func(f Serializer, r *response) ([]byte, error) {
	if f == nil {
		panic("serializer param cannot be nil")
	}
	if r == nil {
		panic("response param cannot be nil")
	}
	// MarshalIndent adds whitespaces to the encoded JSON.
	// No line prefix ("") and two white spaces indents ("\t") for each element.
	js, err := f(r.body, "", "  ")
	if err != nil {
		return nil, err
	}

	// Append a new line making it easier to view in terminal applications.
	js = append(js, '\n')
	return js, nil
}

// Write sends a JSON response to the client.
var write = func(w http.ResponseWriter, f Serializer, r *response) error {
	if w == nil {
		panic("response writer param cannot be nil")
	}
	if r == nil {
		panic("response param cannot be nil")
	}

	js, err := serialize(f, r)
	if err != nil {
		return err
	}

	// At this point it is safe to add any headers as we know that we will not
	// encounter any more errors before writing the response.
	// Custom headers value pass by param method.
	for key, value := range r.header {
		w.Header()[key] = value
	}
	// Response Content Type
	w.Header().Set("Content-Type", "application/json")
	// Response Status Code
	w.WriteHeader(r.statusCode)

	w.Write(js)
	return nil
}

var writeError = func(w http.ResponseWriter, f Serializer) Response {
	rw := new(http.StatusInternalServerError, err, "error", internalServerErrorMsg)
	if err := write(w, f, rw); err != nil {
		w.WriteHeader(500)
	}
	return rw
}
