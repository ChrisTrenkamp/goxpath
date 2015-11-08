package test

import (
	"testing"

	"github.com/ChrisTrenkamp/goxpath/eval"
)

func TestPath1(t *testing.T) {
	testXpath := `/test/path`
	testXML := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	res, err := eval.FromStr(testXpath, testXML)
	if err != nil {
		t.Error(err)
		return
	}
	if len(res) != 1 {
		t.Error("Result not 1")
		return
	}
	if res[0].XPathString() != "<path></path>" {
		t.Error("Incorrect result")
	}
}
