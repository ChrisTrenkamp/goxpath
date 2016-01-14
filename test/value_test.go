package test

import (
	"bytes"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/parser"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

func execVal(xp, x string, exp []string, ns map[string]string, t *testing.T) {
	res := xmltree.Exec(parser.MustParse(xp), xmltree.MustParseXML(bytes.NewBufferString(x)), ns)

	if len(res) != len(exp) {
		t.Error("Result length not valid.  Recieved:")
		for i := range res {
			t.Error(xmltree.MarshalStr(res[i]))
		}
		return
	}

	for i := range res {
		r := res[i].String()
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

func TestNodeVal(t *testing.T) {
	p := `/test`
	x := `<?xml version="1.0" encoding="UTF-8"?><test>test<path>path</path>test2</test>`
	exp := []string{"testpathtest2"}
	execVal(p, x, exp, nil, t)
}

func TestAttrVal(t *testing.T) {
	p := `/p1/@test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo" foo="test"><p2/></p1>`
	exp := []string{"foo"}
	execVal(p, x, exp, nil, t)
}

func TestCommentVal(t *testing.T) {
	p := `//comment()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><!-- comment --></p1>`
	exp := []string{` comment `}
	execVal(p, x, exp, nil, t)
}

func TestProcInstVal(t *testing.T) {
	p := `//processing-instruction()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc test?></p1>`
	exp := []string{`test`}
	execVal(p, x, exp, nil, t)
}

func TestNodeNamespaceVal(t *testing.T) {
	p := `/test:p1/namespace::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:test="http://test"/>`
	exp := []string{`http://test`}
	execVal(p, x, exp, nil, t)
}
