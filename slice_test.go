package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSliceContains_Success(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, "4")

	//verify
	require.Equal(t, true, r)
}

func TestSliceContains_False(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, "5")

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_TypeMismatch(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, 4)

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_Nil(t *testing.T) {

	//setup
	slice := []string{"1", "2", "3", "4"}

	//execute
	r := SliceContains(slice, nil)

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_NotASlice(t *testing.T) {

	//setup

	//execute
	r := SliceContains("1", "1")

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_Empty(t *testing.T) {

	//setup
	slice := []string{}

	//execute
	r := SliceContains(slice, "1")

	//verify
	require.Equal(t, false, r)
}

func TestSliceContains_NilSlice(t *testing.T) {

	//setup

	//execute
	r := SliceContains(nil, "1")

	//verify
	require.Equal(t, false, r)
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
