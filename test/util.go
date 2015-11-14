package test

import (
	"testing"

	"github.com/ChrisTrenkamp/goxpath/xpath"
)

func exec(xp, x string, exp []string, t *testing.T) {
	res, err := xpath.FromStr(xp, x)
	if err != nil {
		t.Error(err)
		return
	}

	if len(res) != len(exp) {
		t.Error("Result length not valid.  Recieved:")
		for i := range res {
			t.Error(xpath.Print(res[i]))
		}
		return
	}

	for i := range exp {
		r, err := xpath.Print(res[i])
		if err != nil {
			t.Error(err.Error())
			return
		}
		if r != exp[i] {
			t.Error("Incorrect result:" + r + "\nExpected: " + exp[i])
			return
		}
	}
}
