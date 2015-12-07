package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringPrefixMatch_Match(t *testing.T) {

	//setup
	str := "MyString"
	prefixes := []string{"One", "1", "My", "Other"}

	//execute
	r := StringPrefixMatch(str, prefixes)

	//verify
	require.Equal(t, true, r)
}

func TestStringPrefixMatch_NoMatch(t *testing.T) {

	//setup
	str := "MyString"
	prefixes := []string{"One", "1", "Other"}

	//execute
	r := StringPrefixMatch(str, prefixes)

	//verify
	require.Equal(t, false, r)
}

func TestStringPrefixMatch_Empty(t *testing.T) {

	//setup
	str := ""
	prefixes := []string{"One", "1", "Other", ""}

	//execute
	r := StringPrefixMatch(str, prefixes)

	//verify
	require.Equal(t, true, r)
}

func TestStringPrefixMatch_EmptyPrefixed(t *testing.T) {

	//setup
	str := ""
	prefixes := []string{}

	//execute
	r := StringPrefixMatch(str, prefixes)

	//verify
	require.Equal(t, false, r)
}

func TestRandomStringSelection_Success(t *testing.T) {

	//setup
	l := []rune("ab")

	//execute
	r := RandomStringSelection(5, l)

	//verify
	require.Equal(t, len(r), 5)
}

func TestRandomStringSelection_Letters(t *testing.T) {

	//setup
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	//execute

	//verify
	require.Equal(t, letters, Letters)
}

func TestRandomString_Success(t *testing.T) {

	//setup

	//execute
	r := RandomString(5)

	//verify
	require.Equal(t, len(r), 5)
}
