package test

import (
	"testing"

	"github.com/ChrisTrenkamp/goxpath/xpath"
)

func exec(xp, x string, exp []string, ns map[string]string, t *testing.T) {
	res, err := xpath.FromStr(xp, x, ns)
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

	for i := range res {
		r, err := xpath.Print(res[i])
		if err != nil {
			t.Error(err.Error())
			return
		}
		valid := false
		for j := range exp {
			if r == exp[j] {
				valid = true
			}
		}
		if !valid {
			t.Error("Incorrect result:" + r)
			return
		}
	}
}
