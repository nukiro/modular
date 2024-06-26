package request

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func integer[T int | int8 | int16 | int32 | int64](p Param, size int) (T, error) {
	if p == "" {
		return T(0), fmt.Errorf("empty parameter")
	}

	i, err := strconv.ParseInt(string(p), 10, size)
	if err != nil {
		return T(0), fmt.Errorf("parameter must be a valid number")
	}

	return T(i), nil
}

type Param string

func (p Param) Int64() (int64, error) {
	return integer[int64](p, 64)
}

func (p Param) Int32() (int32, error) {
	return integer[int32](p, 32)
}

func (p Param) Int16() (int16, error) {
	return integer[int16](p, 16)
}

func (p Param) Int8() (int8, error) {
	return integer[int8](p, 8)
}

func (p Param) Int() (int, error) {
	return integer[int](p, 64)
}

func (p Param) String() (string, error) {
	if p == "" {
		return "", fmt.Errorf("empty parameter")
	}
	return string(p), nil
}

type Params map[string]Param

func get(r *http.Request, name string) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(name)
}

func PathParams(r *http.Request, names ...string) Params {
	params := make(Params)
	for _, name := range names {
		params[name] = Param(get(r, name))
	}
	return params
}

func QueryParams(r *http.Request, names ...string) Params {
	params := make(Params)
	qs := r.URL.Query()
	for _, name := range names {
		if v, ok := qs[name]; ok {
			params[name] = Param(v[0])
		}
	}
	return params
}
