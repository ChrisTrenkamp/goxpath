package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ChrisTrenkamp/goxpath/xpath"
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
		fmt.Println("Specify an XPath expression with one or more files")
	}

	if flag.NArg() == 1 {
		runXPath(flag.Arg(0), os.Stdin, ns, *value)
	}

	for i := 1; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file: %s\n", flag.Arg(i))
			continue
		}

		runXPath(flag.Arg(0), f, ns, *value)
	}
}

func runXPath(x string, r io.Reader, ns namespace, value bool) {
	res, err := xpath.FromReader(x, r, ns)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	for i := range res {
		if value {
			fmt.Print(res[i].GetValue())
		} else {
			str, err := xpath.Print(res[i])
			if err != nil {
				panic(err)
			}
			fmt.Print(str)
		}
	}
}
