package gowork

import (
	"github.com/gorilla/context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestSimpleContext_Success(t *testing.T) {

	//setup
	ctx := NewSimpleContext()
	ctx.Put("Hello", "Hola")

	//execute
	r := ctx.Get("Hello")

	//verify
	require.Equal(t, "Hola", r)
}

func TestSimpleContext_Nil(t *testing.T) {

	//setup
	ctx := NewSimpleContext()

	//execute
	r := ctx.Get("Hello")

	//verify
	require.Nil(t, r)
}


func TestGetString_Success(t *testing.T) {

	//setup
	ctx := NewSimpleContext()
	ctx.Put("Hello", "Hola")

	//execute
	r := ctx.GetString("Hello")

	//verify
	require.Equal(t, "Hola", r)
}

func TestGetString_Nil(t *testing.T) {

	//setup
	ctx := NewSimpleContext()

	//execute
	r := ctx.GetString("Hello")

	//verify
	require.Equal(t, "", r)
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
