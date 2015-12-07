package gowork

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	//setup

	retCode := m.Run()

	//teardown

	os.Exit(retCode)
}
