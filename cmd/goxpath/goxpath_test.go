package main

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestStdinVal(t *testing.T) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	rec = false
	value = true
	ns = make(namespace)
	args = []string{"/root/tag"}
	stdin = strings.NewReader(xml.Header + "<root><tag>test</tag></root>")
	stdout = out
	stderr = err
	exec()
	if out.String() != "test\n" {
		t.Error("Expecting 'test' for the result.  Recieved: ", out.String())
	}
	if retCode != 0 {
		t.Error("Incorrect return value")
	}
}

func TestStdinNonVal(t *testing.T) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	rec = false
	value = false
	ns = make(namespace)
	args = []string{"/root/tag"}
	stdin = strings.NewReader(xml.Header + "<root><tag>test</tag></root>")
	stdout = out
	stderr = err
	exec()
	if out.String() != "<tag>test</tag>\n" {
		t.Error("Expecting '<tag>test</tag>' for the result.  Recieved: ", out.String())
	}
	if retCode != 0 {
		t.Error("Incorrect return value")
	}
}

func TestFile(t *testing.T) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	rec = false
	value = false
	ns = make(namespace)
	ns["foo"] = "http://foo.bar"
	args = []string{"/foo:test/foo:path", "test/1.xml"}
	stdin = strings.NewReader("")
	stdout = out
	stderr = err
	exec()
	if out.String() != `<path xmlns="http://foo.bar">path</path>`+"\n" {
		t.Error(`Expecting '<path xmlns="http://foo.bar">path</path>' for the result.  Recieved: `, out.String())
	}
	if retCode != 0 {
		t.Error("Incorrect return value")
	}
}

func TestDir(t *testing.T) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	rec = true
	value = false
	ns = make(namespace)
	args = []string{"/foo", "test/subdir"}
	stdin = strings.NewReader("")
	stdout = out
	stderr = err
	exec()
	val := strings.Replace(out.String(), "test\\subdir\\", "test/subdir/", -1)
	if val != `test/subdir/2.xml:<foo>bar</foo>`+"\n"+`test/subdir/3.xml:<foo>bar2</foo>`+"\n" {
		t.Error(`Incorrect result.  Recieved: `, val)
	}
	if retCode != 0 {
		t.Error("Incorrect return value")
	}
}

func TestDirNonRec(t *testing.T) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	rec = false
	value = false
	ns = make(namespace)
	args = []string{"/foo", "test/subdir"}
	stdin = strings.NewReader("")
	stdout = out
	stderr = err
	exec()
	val := strings.Replace(err.String(), "test\\subdir\\", "test/subdir/", -1)
	if val != `test/subdir: Is a directory`+"\n" {
		t.Error(`Incorrect result.  Recieved: `, val)
	}
	if retCode != 1 {
		t.Error("Incorrect return value")
	}
}
