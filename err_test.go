package goxpath

import (
	"testing"
)

func TestErr(t *testing.T) {
	xp := `/test/chil::p2`
	_, err := Parse(xp)
	if err.Error() != "Invalid Axis specifier, chil" {
		t.Error("Invalid error message:" + err.Error())
	}
}
