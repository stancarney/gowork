package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNotFoundError_Error_CustomMessage(t *testing.T) {

	//setup
	e := NotFoundError("custom")

	//execute
	str := e.Error()

	//verify
	require.Equal(t, "custom", str)
}

func TestNotFoundError_Error_DefaultMessage(t *testing.T) {

	//setup
	e := NewNotFoundError()

	//execute
	str := e.Error()

	//verify
	require.Equal(t, "not found", str)
}

func TestStaleEntityError_Error_CustomMessage(t *testing.T) {

	//setup
	e := StaleEntityError("custom")


	//execute
	str := e.Error()

	//verify
	require.Equal(t, "custom", str)
}

func TestStaleEntityError_Error_DefaultMessage(t *testing.T) {

	//setup
	e := NewStaleEntityError()

	//execute
	str := e.Error()

	//verify
	require.Equal(t, STALE_ENTITY_MSG, str)
}