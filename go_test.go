package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetFunctionName_Success(t *testing.T) {

	//setup
	f := TestGetFunctionName_Success


	//execute
	r := FunctionName(f)

	//verify
	require.Equal(t, "gowork.TestGetFunctionName_Success", r)
}

func TestGetFunctionName_Pointer(t *testing.T) {

	//setup
	f := TestGetFunctionName_Success


	//execute
	r := FunctionName(&f)

	//verify
	require.Equal(t, "gowork.TestGetFunctionName_Success", r)
}

func TestGetFunctionName_Nil(t *testing.T) {

	//setup


	//execute
	r := FunctionName(nil)

	//verify
	require.Equal(t, "", r)
}

func TestGetCurrentFunctionName_Success(t *testing.T) {

	//setup

	//execute
	r := CurrentFunctionName(2)

	//verify
	require.Equal(t, "gowork.TestGetCurrentFunctionName_Success", r)
}

func TestGetStructName_Success(t *testing.T) {

	//setup
	type blah struct {}
	s := blah{}

	//execute
	r := StructName(s)

	//verify
	require.Equal(t, "gowork.blah", r)
}

func TestGetStructName_Pointer(t *testing.T) {

	//setup
	type blah struct {}
	s := blah{}

	//execute
	r := StructName(&s)

	//verify
	require.Equal(t, "gowork.blah", r)
}

func TestGetStructName_Nil(t *testing.T) {

	//setup

	//execute
	r := StructName(nil)

	//verify
	require.Equal(t, "", r)
}

func TestGetStringValue_Success(t *testing.T) {

	//setup
	type blah struct {
		Value string
	}
	s := blah{"foo"}

	//execute
	r := StringValue(s, "Value")

	//verify
	require.Equal(t, "foo", r)
}

func TestGetStringValue_Pointer(t *testing.T) {

	//setup
	type blah struct {
		Value string
	}
	s := blah{"foo"}

	//execute
	r := StringValue(&s, "Value")

	//verify
	require.Equal(t, "foo", r)
}

func TestGetStringValue_Nil(t *testing.T) {

	//setup

	//execute
	r := StringValue(nil, "Value")

	//verify
	require.Equal(t, "", r)
}


/*
func TestSimpleContext_Nil(t *testing.T) {

	//setup
	ctx := NewSimpleContext()

	//execute
	r := ctx.Get("Hello")

	//verify
	require.Nil(t, r)
}

func TestGetContext_Success(t *testing.T) {

	//setup
	u, _ := url.ParseRequestURI("http://localhost")
	request := &http.Request{URL: u}

	ctx := NewSimpleRequestContext(request)
	context.Set(request, ReqCtx, ctx)

	//execute
	r := GetContext(request)

	//verify
	require.NotNil(t, r)
	src := r.(*SimpleRequestContext)
	require.Equal(t, src.Request, request)
}
*/
