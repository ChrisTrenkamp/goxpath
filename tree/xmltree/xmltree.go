package xmltree

import (
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/net/html/charset"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlbuilder"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlele"
)

//ParseOptions is a set of methods and function pointers that alter
//the way the XML decoder works and the Node types that are created.
//Options that are not set will default to what is set in internal/defoverride.go
type ParseOptions struct {
	Strict  bool
	XMLRoot func() xmlbuilder.XMLBuilder
}

//DirectiveParser is an optional interface extended from XMLBuilder that handles
//XML directives.
type DirectiveParser interface {
	xmlbuilder.XMLBuilder
	Directive(xml.Directive, *xml.Decoder)
}

//ParseSettings is a function for setting the ParseOptions you want when
//parsing an XML tree.
type ParseSettings func(s *ParseOptions)

//MustParseXML is like ParseXML, but panics instead of returning an error.
func MustParseXML(r io.Reader, op ...ParseSettings) tree.Node {
	ret, err := ParseXML(r, op...)

	if err != nil {
		panic(err)
	}

	return ret
}

//ParseXML creates an XMLTree structure from an io.Reader.
func ParseXML(r io.Reader, op ...ParseSettings) (tree.Node, error) {
	ov := ParseOptions{
		Strict:  true,
		XMLRoot: xmlele.Root,
	}
	for _, i := range op {
		i(&ov)
	}

	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = ov.Strict

	ordrPos := 1
	xmlTree := ov.XMLRoot()

	t, err := dec.Token()

	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("Premature end of XML file")
		}
		return nil, err
	}

	brokenHeader := false
	switch t := t.(type) {
	case xml.ProcInst:
		if t.Target != "xml" {
			brokenHeader = true
		}
	default:
		brokenHeader = true
	}

	if brokenHeader {
		if ov.Strict {
			return nil, fmt.Errorf("Malformed XML file")
		}
	} else {
		t, err = dec.Token()
	}

	opts := xmlbuilder.BuilderOpts{
		Dec: dec,
	}

	for err == nil {
		switch xt := t.(type) {
		case xml.StartElement:
			setEle(&opts, xmlTree, xt, &ordrPos)
			if xmlTree.GetNodeType() == tree.NtRoot {
				opts.NS[xml.Name{Space: "", Local: "xml"}] = tree.XMLSpace
				opts.AttrStartPos++
				ordrPos++
			}
			xmlTree = xmlTree.CreateNode(&opts)
		case xml.CharData:
			setNode(&opts, xmlTree, xt, tree.NtChd, &ordrPos)
			xmlTree = xmlTree.CreateNode(&opts)
		case xml.Comment:
			setNode(&opts, xmlTree, xt, tree.NtComm, &ordrPos)
			xmlTree = xmlTree.CreateNode(&opts)
		case xml.ProcInst:
			setNode(&opts, xmlTree, xt, tree.NtPi, &ordrPos)
			xmlTree = xmlTree.CreateNode(&opts)
		case xml.EndElement:
			xmlTree = xmlTree.EndElem()
		case xml.Directive:
			if dp, ok := xmlTree.(DirectiveParser); ok {
				dp.Directive(xt.Copy(), dec)
			}
		}

		t, err = dec.Token()
	}

	if err == io.EOF {
		err = nil
	}

	return xmlTree, err
}

func setEle(opts *xmlbuilder.BuilderOpts, xmlTree xmlbuilder.XMLBuilder, ele xml.StartElement, ordrPos *int) {
	opts.NodePos = *ordrPos
	opts.Tok = ele
	opts.Attrs = opts.Attrs[0:0:cap(opts.Attrs)]
	opts.NS = make(map[xml.Name]string)
	opts.NodeType = tree.NtEle

	*ordrPos++

	for i := range ele.Attr {
		attr := ele.Attr[i].Name
		val := ele.Attr[i].Value

		if (attr.Local == "xmlns" && attr.Space == "") || attr.Space == "xmlns" {
			opts.NS[attr] = val
		} else {
			opts.Attrs = append(opts.Attrs, &ele.Attr[i])
		}
	}

	attrStart := *ordrPos

	if nstree, ok := xmlTree.(tree.NSElem); ok {
		ns := tree.BuildNS(nstree)

		for _, i := range ns {
			if _, ok := opts.NS[i.Attr.Name]; !ok {
				attrStart++
			}
		}

		if _, ok := opts.NS[xml.Name{Space: "", Local: "xmlns"}]; ok {
			attrStart--
		}

		attrStart += len(opts.NS)
	}

	attrStart += len(ele.Attr) - len(opts.NS)
	opts.AttrStartPos = attrStart
	*ordrPos = attrStart + 1
}

func setNode(opts *xmlbuilder.BuilderOpts, xmlTree xmlbuilder.XMLBuilder, tok xml.Token, nt tree.NodeType, ordrPos *int) {
	opts.Tok = xml.CopyToken(tok)
	opts.NodeType = nt
	opts.NodePos = *ordrPos
	*ordrPos++
}
