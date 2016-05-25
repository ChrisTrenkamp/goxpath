package goxpath

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

func execErr(xp, x string, errStr string, ns map[string]string, t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicked: from XPath expr: '" + xp)
			t.Error(r)
			t.Error(string(debug.Stack()))
		}
	}()
	_, err := ExecStr(xp, xmltree.MustParseXML(bytes.NewBufferString(x)), ns)

	if err.Error() != errStr {
		t.Error("Incorrect result:'" + err.Error() + "' from XPath expr: '" + xp + "'.  Expecting: '" + errStr + "'")
		return
	}
}

func TestErr(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`/test/chil::p2`, x, "Invalid Axis specifier, chil", nil, t)
}

func TestUnknownFunction(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`invFunc()`, x, "Unknown function: invFunc", nil, t)
}

func TestUnterminatedString(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`"asdf`, x, "Unexpected end of string literal.", nil, t)
}

func TestUnterminatedParenths(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`(1 + 2 * 3`, x, "Missing end )", nil, t)
}

func TestUnterminatedNTQuotes(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?foo bar?></p1>`
	execErr(`//processing-instruction('foo)`, x, "Unexpected end of string literal.", nil, t)
}

func TestUnterminatedNTParenths(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?foo bar?></p1>`
	execErr(`//processing-instruction('foo'`, x, "Missing ) at end of NodeType declaration.", nil, t)
}

func TestUnterminatedFnParenths(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`true(`, x, "Missing ) at end of function declaration.", nil, t)
}

func TestEmptyPred(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execErr(`/p1[ ]`, x, "Missing content in predicate.", nil, t)
}

func TestUnterminatedPred(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execErr(`/p1[. = 'text'`, x, "Missing ] at end of predicate.", nil, t)
}

func TestNotEnoughArgs(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execErr(`concat('test')`, x, "Invalid number of arguments", nil, t)
}
