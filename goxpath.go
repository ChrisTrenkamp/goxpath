package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ChrisTrenkamp/goxpath/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlres"
)

type namespace map[string]string

func (n *namespace) String() string {
	return fmt.Sprint(*n)
}

func (n *namespace) Set(value string) error {
	nsMap := strings.Split(value, "=")
	if len(nsMap) != 2 {
		return fmt.Errorf("Invalid namespace mapping: " + value)
	}
	(*n)[nsMap[0]] = nsMap[1]
	return nil
}

func main() {
	ns := make(namespace)
	value := flag.Bool("v", false, "Output the string value of the XPath result")

	flag.Var(&ns, "ns", "Namespace mappings. e.g. -ns myns=http://example.com")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Specify an XPath expression with one or more files, or pipe the XML from stdin.")
		os.Exit(1)
	}

	xp, err := goxpath.Parse(flag.Arg(0))

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	if flag.NArg() == 1 {
		ret, err := runXPath(xp, os.Stdin, ns, *value)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		for _, i := range ret {
			fmt.Println(i)
		}
	}

	hasErr := err == nil

	for i := 1; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file: %s\n", flag.Arg(i))
			hasErr = true
			continue
		}

		ret, err := runXPath(xp, f, ns, *value)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			hasErr = true
		}

		for _, j := range ret {
			if len(flag.Args()) > 2 {
				fmt.Printf("%s: ", flag.Arg(i))
			}
			fmt.Println(j)
		}
	}

	if hasErr {
		os.Exit(1)
	}
}

func runXPath(x goxpath.XPathExec, r io.Reader, ns namespace, value bool) ([]string, error) {
	t, err := xmltree.ParseXML(r)

	if err != nil {
		return nil, err
	}

	res, err := goxpath.Exec(x, t, ns)

	if err != nil {
		return nil, err
	}

	ret := make([]string, len(res))

	for i := range res {
		if _, ok := res[i].(xmlres.XMLPrinter); !ok || value {
			ret[i] = res[i].String()
		} else {
			buf := bytes.Buffer{}
			err = xmltree.Marshal(res[i].(xmlres.XMLPrinter), &buf)

			if err != nil {
				return nil, err
			}

			ret[i] = buf.String()
		}
	}

	return ret, nil
}
