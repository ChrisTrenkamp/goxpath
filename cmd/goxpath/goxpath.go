package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree"
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

var rec = flag.Bool("r", false, "Recursive")
var retCode = 0

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
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	if flag.NArg() == 1 {
		ret, err := runXPath(xp, os.Stdin, ns, *value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		for _, i := range ret {
			fmt.Println(i)
		}
	}

	for i := 1; i < flag.NArg(); i++ {
		procPath(flag.Arg(i), xp, ns, *value)
	}

	os.Exit(retCode)
}

func procPath(path string, x goxpath.XPathExec, ns namespace, value bool) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file: %s\n", path)
		retCode = 1
		return
	}

	fi, err := f.Stat()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file: %s\n", path)
		retCode = 1
		return
	}

	if fi.IsDir() {
		procDir(path, x, ns, value)
		return
	}

	ret, err := runXPath(x, f, ns, value)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", path, err.Error())
		retCode = 1
	}

	for _, j := range ret {
		if len(flag.Args()) > 2 || *rec {
			fmt.Printf("%s: ", path)
		}

		fmt.Println(j)
	}
}

func procDir(path string, x goxpath.XPathExec, ns namespace, value bool) {
	if !*rec {
		fmt.Fprintf(os.Stderr, "%s: Is a directory\n", path)
		retCode = 1
		return
	}

	list, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read directory: %s\n", path)
		retCode = 1
		return
	}

	for _, i := range list {
		procPath(filepath.Join(path, i.Name()), x, ns, value)
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
		if _, ok := res[i].(tree.Node); !ok || value {
			ret[i] = res[i].String()
		} else {
			buf := bytes.Buffer{}
			err = goxpath.Marshal(res[i].(tree.Node), &buf)

			if err != nil {
				return nil, err
			}

			ret[i] = buf.String()
		}
	}

	return ret, nil
}
