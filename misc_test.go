package goxpath

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlele"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlnode"
)

func TestISO_8859_1(t *testing.T) {
	p := `/test`
	x := `<?xml version="1.0" encoding="iso-8859-1"?><test>test<path>path</path>test2</test>`
	exp := "testpathtest2"
	execVal(p, x, exp, nil, t)
}

func TestNodePos(t *testing.T) {
	ns := map[string]string{"test": "http://test", "test2": "http://test2", "test3": "http://test3"}
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" attr1="foo"><p2 xmlns="http://test2" xmlns:test="http://test3" attr2="bar">text</p2></p1>`
	testPos := func(path string, pos int) {
		res := MustParse(path).MustExec(xmltree.Adapter{}, xmltree.MustParseXML(bytes.NewBufferString(x)), func(o *Opts) { o.NS = ns }).(tree.NodeSet)
		if len(res.GetNodes()) != 1 {
			t.Errorf("Result length not 1: %s", path)
			return
		}
		exPos := xmltree.Adapter{}.NodePos(res.GetNodes()[0])
		if exPos != pos {
			t.Errorf("Node position not correct.  Recieved %d, expected %d", exPos, pos)
		}
	}
	testPos("/", 0)
	testPos("/*", 1)
	testPos("/*/namespace::*[1]", 2)
	testPos("/*/namespace::*[2]", 3)
	testPos("/*/attribute::*[1]", 4)
	testPos("//*:p2", 5)
	testPos("//*:p2/namespace::*[1]", 6)
	testPos("//*:p2/namespace::*[2]", 7)
	testPos("//*:p2/namespace::*[3]", 8)
	testPos("//*:p2/attribute::*[1]", 9)
	testPos("//text()", 10)
}

func TestNSSort(t *testing.T) {
	testNS := func(n xmlnode.Node, url string) {
		if n.(xmlele.NS).Value != url {
			t.Errorf("Unexpected namespace %s.  Expecting %s", n.(xmlele.NS).Value, url)
		}
	}
	ns := map[string]string{"test": "http://test", "test2": "http://test2", "test3": "http://test3"}
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" xmlns:test2="http://test2" xmlns:test3="http://test3" attr2="bar"/>`
	res := MustParse("/*:p1/namespace::*").MustExec(xmltree.Adapter{}, xmltree.MustParseXML(bytes.NewBufferString(x)), func(o *Opts) { o.NS = ns }).(tree.NodeSet)
	testNS(res.GetNodes()[0].(xmlnode.Node), ns["test"])
	testNS(res.GetNodes()[1].(xmlnode.Node), ns["test2"])
	testNS(res.GetNodes()[2].(xmlnode.Node), ns["test3"])
	testNS(res.GetNodes()[3].(xmlnode.Node), "http://www.w3.org/XML/1998/namespace")
}

func TestFindNodeByPos(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" attr1="foo"><p2 xmlns="http://test2" xmlns:test="http://test3" attr2="bar"><p3/>text<p4/></p2></p1>`
	nt := xmltree.MustParseXML(bytes.NewBufferString(x))
	if xmltree.FindNodeByPos(nt, 5).GetNodeType() != tree.NtElem {
		t.Error("Node 5 not element")
	}
	if xmltree.FindNodeByPos(nt, 14).GetNodeType() != tree.NtChd {
		t.Error("Node 14 not char data")
	}
	if xmltree.FindNodeByPos(nt, 4).GetNodeType() != tree.NtAttr {
		t.Error("Node 4 not attribute")
	}
	if xmltree.FindNodeByPos(nt, 3).GetNodeType() != tree.NtNs {
		t.Error("Node 3 not namespace")
	}
	if xmltree.FindNodeByPos(nt, 19) != nil {
		t.Error("Invalid node returned")
	}
}

//GetAttribute is a convenience function for getting the specified attribute from an element.
//false is returned if the attribute is not found.
func GetAttribute(n xmlnode.Elem, local, space string) (xml.Attr, bool) {
	attrs := n.GetAttrs()
	for _, i := range attrs {
		attr := i.GetToken().(xml.Attr)
		if local == attr.Name.Local && space == attr.Name.Space {
			return attr, true
		}
	}
	return xml.Attr{}, false
}

//GetAttributeVal is like GetAttribute, except it returns the attribute's value.
func GetAttributeVal(n xmlnode.Elem, local, space string) (string, bool) {
	attr, ok := GetAttribute(n, local, space)
	return attr.Value, ok
}

//GetAttrValOrEmpty is like GetAttributeVal, except it returns an empty string if
//the attribute is not found instead of false.
func GetAttrValOrEmpty(n xmlnode.Elem, local, space string) string {
	val, ok := GetAttributeVal(n, local, space)
	if !ok {
		return ""
	}
	return val
}

func TestFindAttr(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:test="http://test" attr1="foo" test:attr2="bar" />`
	nt := xmltree.MustParseXML(bytes.NewBufferString(x))
	res, _ := ParseExec("/p1", xmltree.Adapter{}, nt)
	node := res.(tree.NodeSet).GetNodes()[0].(xmlnode.Elem)
	if val, ok := GetAttributeVal(node, "attr1", ""); !ok || val != "foo" {
		t.Error("attr1 not foo")
	}
	if val, ok := GetAttributeVal(node, "attr2", "http://test"); !ok || val != "bar" {
		t.Error("attr2 not bar")
	}
	if val, ok := GetAttributeVal(node, "attr3", ""); ok || val != "" {
		t.Error("attr3 is set")
	}
	if val := GetAttrValOrEmpty(node, "attr3", ""); val != "" {
		t.Error("attr3 is set")
	}
	if val := GetAttrValOrEmpty(node, "attr1", ""); val != "foo" {
		t.Error("attr1 not foo")
	}
}

func TestVariable(t *testing.T) {
	x := xmltree.MustParseXML(bytes.NewBufferString(xml.Header + "<p1><p2>foo</p2><p3>bar</p3></p1>"))
	xp := MustParse(`/p1/p2`)
	res := xp.MustExec(xmltree.Adapter{}, x)
	opt := func(o *Opts) {
		o.Vars["prev"] = res
	}
	xp = MustParse(`$prev = 'foo'`)
	if res, err := xp.ExecBool(xmltree.Adapter{}, x, opt); err != nil || !res {
		t.Error("Incorrect result", res, err)
	}
	if _, err := xp.ExecBool(xmltree.Adapter{}, x); err == nil {
		t.Error("Error not nil")
	}
	if _, err := Parse(`$ = 'foo'`); err == nil {
		t.Error("Parse error not nil")
	}
}
