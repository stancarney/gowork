package gowork

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	validator "gopkg.in/validator.v2"
	"net/http/httptest"
	"testing"
	"errors"
)

type TestStruct struct {
	StringValue string `validate:"nonzero"`
	IntValue    int    `validate:"nonzero"`
}

type BadStruct struct {
	StringValue string
	IntValue    int
}

func (s BadStruct) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("BOOM!")
}

func (s BadStruct) UnmarshalJSON(b []byte) error {
	return fmt.Errorf("BOOM!")
}

func TestWriteJSON_Valid(t *testing.T) {

	//setup
	w := httptest.NewRecorder()
	data := TestStruct{"Hello", 5}

	//execute
	WriteJSON(w, data)

	//verify
	require.Equal(t, 200, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"StringValue\":\"Hello\",\"IntValue\":5}"), w.Body.String())
}

func TestWriteJSON_Error(t *testing.T) {

	//setup
	w := httptest.NewRecorder()
	data := BadStruct{"Hello", 5}

	//execute
	WriteJSON(w, data)

	//verify
	require.Equal(t, 500, w.Code)
	require.Equal(t, "json: error calling MarshalJSON for type gowork.BadStruct: BOOM!\n", w.Body.String())
}

func TestWriteMessageToJSON_Valid(t *testing.T) {

	//setup
	w := httptest.NewRecorder()

	//execute
	WriteMessageToJSON(w, "Hello")

	//verify
	require.Equal(t, 200, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"msg\":\"Hello\"}"), w.Body.String())
}

func TestWriteErrorToJSON_Default(t *testing.T) {

	//setup
	w := httptest.NewRecorder()
	data := fmt.Errorf("My Error")

	//execute
	WriteErrorToJSON(w, 500, data)

	//verify
	require.Equal(t, 500, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"errors\":[{\"msg\":\"My Error\"}],\"count\":1}"), w.Body.String())
}

func TestWriteErrorToJSON_ErrorMap(t *testing.T) {

	//setup
	w := httptest.NewRecorder()
	testStruct := TestStruct{}

	data := validator.Validate(testStruct)

	//execute
	WriteErrorToJSON(w, 500, data)

	//verify
	require.Equal(t, 500, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])

	//TODO:Stan look at pulling this into a custom assert maybe?
	j := make(map[string]interface{})
	json.Unmarshal(w.Body.Bytes(), &j)
	require.Equal(t, 2, len(j["errors"].([]interface{})))
	require.NotEmpty(t, j["errors"].([]interface{})[0].(map[string]interface{})["field"])
	require.Equal(t, "zero value", j["errors"].([]interface{})[0].(map[string]interface{})["msg"])
	require.NotEmpty(t, j["errors"].([]interface{})[1].(map[string]interface{})["field"])
	require.Equal(t, "zero value", j["errors"].([]interface{})[1].(map[string]interface{})["msg"])
}

func TestWriteErrorToJSON_String(t *testing.T) {

	//setup
	w := httptest.NewRecorder()

	//execute
	WriteErrorToJSON(w, 500, "My Error")

	//verify
	require.Equal(t, 500, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"errors\":[{\"msg\":\"My Error\"}],\"count\":1}"), w.Body.String())
}

func TestWriteErrorToJSON_Slice(t *testing.T) {

	//setup
	w := httptest.NewRecorder()

	errors := []error{errors.New("one"), errors.New("two")}

	//execute
	WriteErrorToJSON(w, 500, errors)

	//verify
	require.Equal(t, 500, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"errors\":[{\"msg\":\"one\"},{\"msg\":\"two\"}],\"count\":2}"), w.Body.String())
}

func TestJSONResponse_NoError(t *testing.T) {

	//setup
	w := httptest.NewRecorder()

	//execute
	JSONResponse(w, "My Data", nil)

	//verify
	require.Equal(t, 200, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, "\"My Data\"", w.Body.String())
}

func TestJSONResponse_NotFoundError(t *testing.T) {

	//setup
	w := httptest.NewRecorder()

	//execute
	JSONResponse(w, "My Data", NewNotFoundError())

	//verify
	require.Equal(t, 404, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"errors\":[{\"msg\":\"not found\"}],\"count\":1}"), w.Body.String())
}

func TestJSONResponse_Error(t *testing.T) {

	//setup
	w := httptest.NewRecorder()

	//execute
	JSONResponse(w, "My Data", errors.New("Blah"))

	//verify
	require.Equal(t, 500, w.Code)
	require.Equal(t, 1, len(w.HeaderMap))
	require.Equal(t, "application/json", w.HeaderMap["Content-Type"][0])
	require.Equal(t, FormatJSON("{\"errors\":[{\"msg\":\"Blah\"}],\"count\":1}"), w.Body.String())
}