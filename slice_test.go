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

func TestStringMapToSlice_Success(t *testing.T) {

	//setup
	m := make(map[string]interface{})
	m["1"] = &Session{Id: "1", Version: 1}
	m["2"] = &Session{Id: "2", Version: 2}

	//execute
	r := StringMapToSlice(m)

	//verify
	sessions := r.([]Session)
	require.Equal(t, 2, len(sessions))
	require.Equal(t, true, sessions[0].Id == "1" || sessions[1].Id == "1")
	require.Equal(t, true, sessions[0].Id == "2" || sessions[1].Id == "2")
}