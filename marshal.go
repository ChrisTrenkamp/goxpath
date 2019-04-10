package goxpath

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/ChrisTrenkamp/goxpath/tree"
)

//Marshal prints the result tree, r, in XML form to w.
func Marshal(a tree.Adapter, n interface{}, w io.Writer) error {
	return marshal(a, n, w)
}

//MarshalStr is like Marhal, but returns a string.
func MarshalStr(a tree.Adapter, n interface{}) (string, error) {
	ret := bytes.NewBufferString("")
	err := marshal(a, n, ret)

	return ret.String(), err
}

func marshal(a tree.Adapter, n interface{}, w io.Writer) error {
	e := xml.NewEncoder(w)
	err := encTok(a, n, e)
	if err != nil {
		return err
	}

	return e.Flush()
}

func encTok(a tree.Adapter, n interface{}, e *xml.Encoder) error {
	switch a.GetNodeType(n) {
	case tree.NtAttr:
		return encAttr(a.GetAttrTok(n), e)
	case tree.NtElem:
		return encEle(a, n, e)
	case tree.NtNs:
		return encNS(a.GetNamespaceTok(n), e)
	case tree.NtRoot:
		a.ForEachChild(n, func(in interface{}) {
			encTok(a, in, e)
		})
		return nil
	case tree.NtChd:
		return e.EncodeToken(a.GetCharDataTok(n))
	case tree.NtComm:
		return e.EncodeToken(a.GetCommentTok(n))
	case tree.NtPi:
		return e.EncodeToken(a.GetProcInstTok(n))
	}
	return nil
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

func encEle(a tree.Adapter, node interface{}, e *xml.Encoder) error {
	startEle := a.GetElemTok(node)
	ele := xml.StartElement{
		Name: startEle.Name,
	}

	a.ForEachAttr(node, func(attr xml.Attr, ptr interface{}) {
		ele.Attr = append(ele.Attr, attr)
	})

	err := e.EncodeToken(ele)
	if err != nil {
		return err
	}

	a.ForEachChild(node, func(in interface{}) {
		encTok(a, in, e)
	})

	return e.EncodeToken(ele.End())
}
