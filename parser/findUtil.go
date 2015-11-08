package parser

import "encoding/xml"

func (x *xmlTree) findTag(p pathExpr) []*xmlTree {
	ret := []*xmlTree{}

	for i := range x.children {
		switch x.children[i].value.(type) {
		case xml.StartElement:
			se := x.children[i].value.(xml.StartElement)

			if se.Name == p.name {
				ret = append(ret, x.children[i])
			}

			if p.abbr {
				ret = append(ret, x.children[i].findTag(p)...)
			}
		}
	}

	return ret
}
