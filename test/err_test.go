package test

import (
	"testing"

	"github.com/ChrisTrenkamp/goxpath/goxpath"
)

func TestErr(t *testing.T) {
	xp := `/test/chil::p2`
	_, err := goxpath.Parse(xp)
	if err.Error() != "Invalid Axis specifier, chil" {
		t.Error("Invalid error message:" + err.Error())
	}
}
