package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	"github.com/ChrisTrenkamp/goxpath/xtypes"
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

var rec bool
var value bool
var ns = make(namespace)
var unstrict bool
var noFileName bool
var args = []string{}
var stdin io.Reader = os.Stdin
var stdout io.ReadWriter = os.Stdout
var stderr io.ReadWriter = os.Stderr

var retCode = 0

func init() {
	flag.BoolVar(&rec, "r", false, "Recursive")
	flag.BoolVar(&value, "v", false, "Output the string value of the XPath result")
	flag.Var(&ns, "ns", "Namespace mappings. e.g. -ns myns=http://example.com")
	flag.BoolVar(&unstrict, "u", false, "Turns off strict XML validation")
	flag.BoolVar(&noFileName, "h", false, "Suppress filename prefixes.")
}

func main() {
	flag.Parse()
	args = flag.Args()
	exec()
	os.Exit(retCode)
}

func exec() {
	if len(args) < 1 {
		fmt.Fprintf(stdout, "Specify an XPath expression with one or more files, or pipe the XML from stdin.\n")
		retCode = 1
		return
	}

	xp, err := goxpath.Parse(args[0])

	if err != nil {
		fmt.Fprintf(stderr, "%s\n", err.Error())
		retCode = 1
		return
	}

	if len(args) == 1 {
		ret, err := runXPath(xp, stdin, ns, value)
		if err != nil {
			fmt.Fprintf(stderr, "%s\n", err.Error())
			retCode = 1
		}

		printResult(ret, "")
	}

	for i := 1; i < len(args); i++ {
		procPath(args[i], xp, ns, value)
	}
}

func procPath(path string, x goxpath.XPathExec, ns namespace, value bool) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Fprintf(stderr, "Could not open file: %s\n", path)
		retCode = 1
		return
	}

	fi, err := f.Stat()

	if err != nil {
		fmt.Fprintf(stderr, "Could not open file: %s\n", path)
		retCode = 1
		return
	}

	if fi.IsDir() {
		procDir(path, x, ns, value)
		return
	}

	ret, err := runXPath(x, f, ns, value)

	if err != nil {
		fmt.Fprintf(stderr, "%s: %s\n", path, err.Error())
		retCode = 1
	}

	printResult(ret, path)
}

func printResult(ret []string, path string) {
	for _, j := range ret {
		if (len(flag.Args()) > 2 || rec) && !noFileName {
			fmt.Fprintf(stdout, "%s:", path)
		}

		fmt.Fprintf(stdout, "%s\n", j)
	}
}

func procDir(path string, x goxpath.XPathExec, ns namespace, value bool) {
	if !rec {
		fmt.Fprintf(stderr, "%s: Is a directory\n", path)
		retCode = 1
		return
	}

	list, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Fprintf(stderr, "Could not read directory: %s\n", path)
		retCode = 1
		return
	}

	for _, i := range list {
		procPath(filepath.Join(path, i.Name()), x, ns, value)
	}
}

func runXPath(x goxpath.XPathExec, r io.Reader, ns namespace, value bool) ([]string, error) {
	t, err := xmltree.ParseXML(r, func(o *xmltree.ParseOptions) {
		o.Strict = !unstrict
	})

	if err != nil {
		return nil, err
	}

	res, err := goxpath.Exec(x, t, ns)

	if err != nil {
		return nil, err
	}

	var ret []string

	if nodes, ok := res.(xtypes.NodeSet); ok && !value {
		ret = make([]string, len(nodes))
		for i, v := range nodes {
			ret[i], err = goxpath.MarshalStr(v)
			ret[i] = strings.Replace(ret[i], "\n", "&#10;", -1)
			if err != nil {
				return nil, err
			}
		}
	} else {
		ret = strings.Split(res.String(), "\n")
	}

	return ret, nil
}
