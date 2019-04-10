package xmltree

import (
	"encoding/xml"
	"io"
	"sort"

	"golang.org/x/net/html/charset"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlbuilder"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlele"
	"github.com/ChrisTrenkamp/goxpath/treeimpl/xmltree/xmlnode"
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
func MustParseXML(r io.Reader, op ...ParseSettings) xmlnode.Node {
	ret, err := ParseXML(r, op...)

	if err != nil {
		panic(err)
	}

	return ret
}

//ParseXML creates an XMLTree structure from an io.Reader.
func ParseXML(r io.Reader, op ...ParseSettings) (xmlnode.Node, error) {
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
		return nil, err
	}

	if head, ok := t.(xml.ProcInst); ok && head.Target == "xml" {
		t, err = dec.Token()
	}

	opts := xmlbuilder.BuilderOpts{
		Dec: dec,
	}

	for err == nil {
		switch xt := t.(type) {
		case xml.StartElement:
			setEle(&opts, xmlTree, xt, &ordrPos)
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
	opts.NodeType = tree.NtElem

	for i := range ele.Attr {
		attr := ele.Attr[i].Name
		val := ele.Attr[i].Value

		if (attr.Local == "xmlns" && attr.Space == "") || attr.Space == "xmlns" {
			opts.NS[attr] = val
		} else {
			opts.Attrs = append(opts.Attrs, &ele.Attr[i])
		}
	}

	if nstree, ok := xmlTree.(xmlnode.NSElem); ok {
		ns := make(map[xml.Name]string)

		for _, i := range xmlele.BuildNS(nstree) {
			ns[i.Name] = i.Value
		}

		for k, v := range opts.NS {
			ns[k] = v
		}

		if ns[xml.Name{Local: "xmlns"}] == "" {
			delete(ns, xml.Name{Local: "xmlns"})
		}

		for k, v := range ns {
			opts.NS[k] = v
		}

		if xmlTree.GetNodeType() == tree.NtRoot {
			opts.NS[xml.Name{Space: "xmlns", Local: "xml"}] = tree.XMLSpace
		}
	}

	opts.AttrStartPos = len(opts.NS) + len(opts.Attrs) + *ordrPos
	*ordrPos = opts.AttrStartPos + 1
}

func setNode(opts *xmlbuilder.BuilderOpts, xmlTree xmlbuilder.XMLBuilder, tok xml.Token, nt tree.NodeType, ordrPos *int) {
	opts.Tok = xml.CopyToken(tok)
	opts.NodeType = nt
	opts.NodePos = *ordrPos
	*ordrPos++
}

//FindNodeByPos finds a node from the given position.  Returns nil if the node
//is not found.
func FindNodeByPos(n xmlnode.Node, pos int) xmlnode.Node {
	if n.Pos() == pos {
		return n
	}

	if elem, ok := n.(xmlnode.Elem); ok {
		chldrn := elem.GetChildren()
		for i := 1; i < len(chldrn); i++ {
			if chldrn[i-1].Pos() <= pos && chldrn[i].Pos() > pos {
				return FindNodeByPos(chldrn[i-1], pos)
			}
		}

		if len(chldrn) > 0 {
			if chldrn[len(chldrn)-1].Pos() <= pos {
				return FindNodeByPos(chldrn[len(chldrn)-1], pos)
			}
		}

		attrs := elem.GetAttrs()
		for _, i := range attrs {
			if i.Pos() == pos {
				return i
			}
		}

		ns := xmlele.BuildNS(elem)
		for _, i := range ns {
			if i.Position == pos {
				return i
			}
		}
	}

	return nil
}

type Adapter struct{}

func (a Adapter) GetNodeType(in interface{}) tree.NodeType {
	return in.(xmlnode.Node).GetNodeType()
}

func (a Adapter) GetParent(in interface{}) interface{} {
	return in.(xmlnode.Node).GetParent()
}

func (a Adapter) GetAttrTok(in interface{}) xml.Attr {
	return in.(xmlnode.Node).GetToken().(xml.Attr)
}

func (a Adapter) GetNamespaceTok(in interface{}) xml.Attr {
	return in.(xmlnode.Node).GetToken().(xml.Attr)
}

func (a Adapter) GetElemTok(in interface{}) xml.StartElement {
	return in.(*xmlele.XMLEle).GetToken().(xml.StartElement)
}

func (a Adapter) GetElementName(in interface{}) xml.Name {
	return in.(*xmlele.XMLEle).GetToken().(xml.StartElement).Name
}

func (a Adapter) GetProcInstTok(in interface{}) xml.ProcInst {
	return in.(xmlnode.Node).GetToken().(xml.ProcInst)
}

func (a Adapter) GetCharDataTok(in interface{}) xml.CharData {
	return in.(xmlnode.Node).GetToken().(xml.CharData)
}

func (a Adapter) GetCommentTok(in interface{}) xml.Comment {
	return in.(xmlnode.Node).GetToken().(xml.Comment)
}

func (a Adapter) NodePos(in interface{}) int {
	return in.(xmlnode.Node).Pos()
}

func (a Adapter) ForEachAttr(in interface{}, f func(xml.Attr, interface{})) {
	if ele, ok := in.(*xmlele.XMLEle); ok {
		for _, attr := range ele.Attrs {
			tok := a.GetAttrTok(attr)
			f(tok, attr)
		}
	}
}

func (a Adapter) ForEachChild(in interface{}, g func(interface{})) {
	if ele, ok := in.(*xmlele.XMLEle); ok {
		for _, ch := range ele.Children {
			g(ch)
		}
	}
}

func (a Adapter) StringValue(in interface{}) string {
	return in.(xmlnode.Node).ResValue()
}

func (a Adapter) GetNamespaces(in interface{}) []interface{} {
	if ele, ok := in.(*xmlele.XMLEle); ok {
		ns := xmlele.BuildNS(ele)
		ret := make([]interface{}, 0, len(ns))
		for _, x := range ns {
			ret = append(ret, x)
		}
		return ret
	}
	return nil
}

type _nodeset struct {
	nodes   []interface{}
	adapter Adapter
}

func (a Adapter) NewNodeSet(nodes []interface{}) tree.NodeSet {
	return &_nodeset{nodes: nodes, adapter: a}
}

func (n _nodeset) GetNodes() []interface{} {
	return n.nodes
}

//String satisfies the Res interface for NodeSet
func (n _nodeset) String() string {
	if len(n.nodes) == 0 {
		return ""
	}

	return n.adapter.StringValue(n.nodes[0])
}

//Bool satisfies the HasBool interface for node-set's
func (n _nodeset) Bool() tree.Bool {
	return tree.Bool(len(n.nodes) > 0)
}

//Num satisfies the HasNum interface for NodeSet's
func (n _nodeset) Num() tree.Num {
	return tree.String(n.String()).Num()
}

func (n _nodeset) Sort() {
	sort.Slice(n.nodes, func(i, j int) bool {
		return n.adapter.NodePos(n.nodes[i]) < n.adapter.NodePos(n.nodes[j])
	})
}

func (n *_nodeset) Unique() {
	x := map[int]interface{}{}
	for _, r := range n.nodes {
		x[n.adapter.NodePos(r)] = r
	}
	out := make([]interface{}, 0, len(x))
	for _, i := range x {
		out = append(out, i)
	}
	n.nodes = out
}
