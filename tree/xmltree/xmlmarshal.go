package xmltree

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlres"
)

//Marshal prints the result tree, r, in XML form to w.
func Marshal(r xmlres.XMLPrinter, w io.Writer) error {
	return marshal(r, w)
}

//MarshalStr is like Marhal, but returns a string.
func MarshalStr(r xmlres.XMLPrinter) (string, error) {
	ret := bytes.NewBufferString("")
	err := marshal(r, ret)

	return ret.String(), err
}

func marshal(r xmlres.XMLPrinter, w io.Writer) error {
	e := xml.NewEncoder(w)
	err := r.XMLPrint(e)
	if err != nil {
		return err
	}

	err = e.Flush()
	if err != nil {
		return err
	}

	return nil
}
