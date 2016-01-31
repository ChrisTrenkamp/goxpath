package test

import "testing"

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
}

func TestLocalName1(t *testing.T) {
	p := `local-name( / * )`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"test"}
	execVal(p, x, exp, nil, t)
	p = `/test/path/ local-name ( ) `
	exp = []string{"path"}
	execVal(p, x, exp, nil, t)
}
