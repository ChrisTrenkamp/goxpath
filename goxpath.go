package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ChrisTrenkamp/goxpath/eval"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Specify an XPath expression with one or more files")
	}

	if flag.NArg() == 1 {
		runXPath(flag.Arg(0), os.Stdin)
	}

	for i := 1; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file: %s\n", flag.Arg(i))
			continue
		}

		runXPath(flag.Arg(0), f)
	}
}

func runXPath(xpath string, r io.Reader) {
	res, err := eval.FromReader(xpath, r)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	for i := range res {
		fmt.Println(res[i].XPathString())
	}
}
