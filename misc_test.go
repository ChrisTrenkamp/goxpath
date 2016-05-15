package goxpath

import (
	"bytes"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
)

func TestNodePos(t *testing.T) {
	ns := map[string]string{"test": "http://test", "test2": "http://test2", "test3": "http://test3"}
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" attr1="foo"><p2 xmlns="http://test2" xmlns:test="http://test3" attr2="bar">text</p2></p1>`
	testPos := func(path string, pos int) {
		res := MustExec(MustParse(path), xmltree.MustParseXML(bytes.NewBufferString(x)), ns).(xtypes.NodeSet)
		if len(res) != 1 {
			t.Errorf("Result length not 1: %s", path)
			return
		}
		exPos := res[0].(tree.Node).Pos()
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
	testNS := func(n tree.Node, url string) {
		if n.(tree.NS).Value != url {
			t.Errorf("Unexpected namespace %s.  Expecting %s", n.(tree.NS).Value, url)
		}
	}
	ns := map[string]string{"test": "http://test", "test2": "http://test2", "test3": "http://test3"}
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" xmlns:test2="http://test2" xmlns:test3="http://test3" attr2="bar"/>`
	res := MustExec(MustParse("/*:p1/namespace::*"), xmltree.MustParseXML(bytes.NewBufferString(x)), ns).(xtypes.NodeSet)
	testNS(res[0], ns["test"])
	testNS(res[1], ns["test2"])
	testNS(res[2], ns["test3"])
	testNS(res[3], "http://www.w3.org/XML/1998/namespace")
}
