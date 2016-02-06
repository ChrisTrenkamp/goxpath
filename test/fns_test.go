package test

import (
	"fmt"
	"testing"
)

func TestStrLit(t *testing.T) {
	p := `'strlit'`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"strlit"}
	execVal(p, x, exp, nil, t)
}

func TestNumLit(t *testing.T) {
	p := `123`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"123"}
	execVal(p, x, exp, nil, t)
	p = `123.456`
	exp = []string{"123.456"}
	execVal(p, x, exp, nil, t)
}

func TestLast(t *testing.T) {
	p := `/p1/*/last()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p3/><p4/></p1>`
	exp := []string{"3"}
	execVal(p, x, exp, nil, t)
	p = `/p1/p5/last()`
	exp = []string{"0"}
	execVal(p, x, exp, nil, t)
}

func TestCount(t *testing.T) {
	p := `count(/p1)`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><?test?></p2><p3/><p4/></p1>`
	exp := []string{"5"}
	execVal(p, x, exp, nil, t)
}

func TestNames(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><test xmlns="http://foo.com" xmlns:bar="http://bar.com" bar:attr="val"><?pi pival?><!--comment--></test>`
	testMap := make(map[string]map[string]string)
	testMap["/*"] = make(map[string]string)
	testMap["/*"]["local-name"] = "test"
	testMap["/*"]["namespace-uri"] = "http://foo.com"
	testMap["/*"]["name"] = "{http://foo.com}test"

	testMap["/*/@*:attr"] = make(map[string]string)
	testMap["/*/@*:attr"]["local-name"] = "attr"
	testMap["/*/@*:attr"]["namespace-uri"] = "http://bar.com"
	testMap["/*/@*:attr"]["name"] = "{http://bar.com}attr"

	testMap["//processing-instruction()"] = make(map[string]string)
	testMap["//processing-instruction()"]["local-name"] = "pi"
	testMap["//processing-instruction()"]["namespace-uri"] = ""
	testMap["//processing-instruction()"]["name"] = "pi"

	testMap["//comment()"] = make(map[string]string)
	testMap["//comment()"]["local-name"] = ""
	testMap["//comment()"]["namespace-uri"] = ""
	testMap["//comment()"]["name"] = ""

	for path, i := range testMap {
		for nt, res := range i {
			p := fmt.Sprintf("%s(%s)", nt, path)
			exp := []string{res}
			execVal(p, x, exp, nil, t)
		}
	}
}

func TestBoolean(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p3/><p4/></p1>`
	execVal(`true()`, x, []string{"true"}, nil, t)
	execVal(`false()`, x, []string{"false"}, nil, t)
	p := `boolean(/p1/p2)`
	exp := []string{"true"}
	execVal(p, x, exp, nil, t)
	p = `boolean(/p1/p5)`
	exp = []string{"false"}
	execVal(p, x, exp, nil, t)
	p = `boolean('123')`
	exp = []string{"true"}
	execVal(p, x, exp, nil, t)
	p = `boolean(123)`
	exp = []string{"true"}
	execVal(p, x, exp, nil, t)
	p = `boolean('')`
	exp = []string{"false"}
	execVal(p, x, exp, nil, t)
	p = `boolean(0)`
	exp = []string{"false"}
	execVal(p, x, exp, nil, t)
}

func TestNot(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p3/><p4/></p1>`
	execVal(`not(false())`, x, []string{"true"}, nil, t)
	execVal(`not(true())`, x, []string{"false"}, nil, t)
}
