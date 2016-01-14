package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ChrisTrenkamp/goxpath/parser"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
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

	xp, err := parser.Parse(flag.Arg(0))

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	if flag.NArg() == 1 {
		runXPath(xp, os.Stdin, ns, *value)
	}

	for i := 1; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file: %s\n", flag.Arg(i))
			continue
		}

		err = runXPath(xp, f, ns, *value)

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func runXPath(x parser.XPathExec, r io.Reader, ns namespace, value bool) error {
	t, err := xmltree.ParseXML(r)

	if err != nil {
		return err
	}

	res := xmltree.Exec(x, t, ns)

	for i := range res {
		if value {
			fmt.Print(res[i])
		} else {
			err = xmltree.Marshal(res[i], os.Stdout)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
