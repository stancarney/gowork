package gowork

import (
	"encoding/json"
	"fmt"
	validator "gopkg.in/validator.v2"
	"net/http"
	"reflect"
	"bytes"
)

type Message struct {
	Message string `json:"msg"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
	Count  int `json:"count"`
}

func (er * ErrorResponse) Add(error Error) {
	er.Errors = append(er.Errors, error)
	er.Count++
}

type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"msg"`
	Index   int `json:"index,omitempty"`
}

type IndexedError struct {
	Err   error
	Index int
}

func (e IndexedError) Error() string {
	return fmt.Sprintf("%d: %s", e.Index, e.Err.Error())
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	str := string(b)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, str)

	return
}

func WriteMessageToJSON(w http.ResponseWriter, message string) {
	data := &Message{Message: message}
	WriteJSON(w, data)
	return
}

func WriteErrorToJSON(w http.ResponseWriter, code int, err interface{}) {

	resp := &ErrorResponse{}

	errorType := reflect.ValueOf(err)
	switch errorType.Kind() {
	case reflect.Slice:
		m := err.([]error)
		for _, v := range m {
			formatError(v, resp)
		}

	default:
		formatError(err, resp)
	}

	w.WriteHeader(code)
	WriteJSON(w, resp)

	return
}

func JSONResponse(w http.ResponseWriter, data interface{}, err interface{}) {

	if _, ok := err.(NotFoundError); ok {
		WriteErrorToJSON(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		WriteErrorToJSON(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSON(w, data)
}

func formatError(err interface{}, resp *ErrorResponse) {

	index := 0
	if e, ok := err.(IndexedError); ok {
		index = e.Index
		err = e.Err
	}

	if e, ok := err.(validator.ErrorMap); ok {
		for k, v := range e {
			resp.Add(Error{Field: k, Message: v.Error(), Index: index})
		}
		return
	}

	if e, ok := err.(error); ok {
		resp.Add(Error{Message: e.Error(), Index: index})
		return
	}

	resp.Add(Error{Message: fmt.Sprintf("%s", err), Index: index})
}

//FormatJSON is a quick utility method to format JSON. Primarily used in logging/testing as it doesn't return an error on failure.
func FormatJSON(src string) string {
	dst := &bytes.Buffer{}
	s := []byte(src)
	if err := json.Indent(dst, s, "", "  "); err != nil {
		return err.Error()
	}
	return dst.String()
}