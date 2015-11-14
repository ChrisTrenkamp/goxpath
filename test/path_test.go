package test

import "testing"

func TestAbsPath(t *testing.T) {
	p := `/test/path`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<path></path>"}
	exec(p, x, exp, t)
}

func TestRelPath(t *testing.T) {
	p := `//path`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<path></path>"}
	exec(p, x, exp, t)
}

func TestParent(t *testing.T) {
	p := `/test/path/parent::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<test><path></path></test>"}
	exec(p, x, exp, t)
}

func TestAncestor(t *testing.T) {
	p := `/p1/p2/p3/p1/ancestor::p1`
	x := `
<?xml version="1.0" encoding="UTF-8"?>
<p1>
	<p2>
		<p3>
			<p1></p1>
		</p3>
	</p2>
</p1>`
	exp := []string{`<p1>
&#x9;<p2>
&#x9;&#x9;<p3>
&#x9;&#x9;&#x9;<p1></p1>
&#x9;&#x9;</p3>
&#x9;</p2>
</p1>`}
	exec(p, x, exp, t)
}

func TestAncestorOrSelf(t *testing.T) {
	p := `/p1/p2/p3/p1/ancestor-or-self::p1`
	x := `
<?xml version="1.0" encoding="UTF-8"?>
<p1>
	<p2>
		<p3>
			<p1></p1>
		</p3>
	</p2>
</p1>`
	exp := []string{`<p1></p1>`, `<p1>
&#x9;<p2>
&#x9;&#x9;<p3>
&#x9;&#x9;&#x9;<p1></p1>
&#x9;&#x9;</p3>
&#x9;</p2>
</p1>`}
	exec(p, x, exp, t)
}

func TestDescendant(t *testing.T) {
	p := `/p1/descendant::p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p1/></p2></p1>`
	exp := []string{`<p1></p1>`}
	exec(p, x, exp, t)
}

func TestDescendantOrSelf(t *testing.T) {
	p := `/p1/descendant-or-self::p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p1/></p2></p1>`
	exp := []string{`<p1><p2><p1></p1></p2></p1>`, `<p1></p1>`}
	exec(p, x, exp, t)
}

func TestAttribute(t *testing.T) {
	p := `/p1/attribute::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"></p1>`
	exp := []string{`<?attribute test="foo"?>`}
	exec(p, x, exp, t)
}

func TestAttributeAbbr(t *testing.T) {
	p := `/p1/@test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"></p1>`
	exp := []string{`<?attribute test="foo"?>`}
	exec(p, x, exp, t)
}

func TestNodeTypeNode(t *testing.T) {
	p := `/p1/child::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p2></p2>`}
	exec(p, x, exp, t)
}

func TestNodeTypeNodeAbbr(t *testing.T) {
	p := `/p1/.`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	exec(p, x, exp, t)
}

func TestNodeTypeParent(t *testing.T) {
	p := `/p1/p2/parent::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	exec(p, x, exp, t)
}

func TestNodeTypeParentAbbr(t *testing.T) {
	p := `/p1/p2/..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	exec(p, x, exp, t)
}
