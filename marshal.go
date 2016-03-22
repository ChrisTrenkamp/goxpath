package goxpath

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Marshal prints the result tree, r, in XML form to w.
func Marshal(n tree.Node, w io.Writer) error {
	return marshal(n, w)
}

//MarshalStr is like Marhal, but returns a string.
func MarshalStr(n tree.Node) (string, error) {
	ret := bytes.NewBufferString("")
	err := marshal(n, ret)

	return ret.String(), err
}

func marshal(n tree.Node, w io.Writer) error {
	e := xml.NewEncoder(w)
	err := encTok(n, e)
	if err != nil {
		return err
	}

	err = e.Flush()
	if err != nil {
		return err
	}

	return nil
}

func encTok(n tree.Node, e *xml.Encoder) error {
	tok := n.GetToken()

	switch n.GetNodeType() {
	case tree.NtAttr:
		return encAttr(tok.(xml.Attr), e)
	case tree.NtEle:
		return encEle(n.(tree.Elem), e)
	case tree.NtNs:
		return encNS(tok.(xml.Attr), e)
	case tree.NtRoot:
		for _, i := range n.(tree.Elem).GetChildren() {
			err := encTok(i, e)
			if err != nil {
				return err
			}
		}
		return nil
	}

	//case tree.NtChd, tree.NtComm, tree.NtPi:
	return e.EncodeToken(tok)
}

func encAttr(a xml.Attr, e *xml.Encoder) error {
	str := a.Name.Local + `="` + a.Value + `"`

	if a.Name.Space != "" {
		str += ` xmlns="` + a.Name.Space + `"`
	}

	pi := xml.ProcInst{
		Target: "attribute",
		Inst:   ([]byte)(str),
	}

	return e.EncodeToken(pi)
}

func encNS(ns xml.Attr, e *xml.Encoder) error {
	pi := xml.ProcInst{
		Target: "namespace",
		Inst:   ([]byte)(ns.Value),
	}
	return e.EncodeToken(pi)
}

func encEle(n tree.Elem, e *xml.Encoder) error {
	val := n.GetToken().(xml.StartElement)

	for i := 0; i < len(val.Attr); i++ {
		if val.Attr[i].Name.Local == "xmlns" || val.Attr[i].Name.Space == "xmlns" {
			val.Attr = append(val.Attr[:i], val.Attr[i+1:]...)
			i--
		}
	}

	err := e.EncodeToken(val)
	if err != nil {
		return err
	}

	if x, ok := n.(tree.Elem); ok {
		for _, i := range x.GetChildren() {
			err := encTok(i, e)
			if err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(xml.EndElement{Name: val.Name})
}
