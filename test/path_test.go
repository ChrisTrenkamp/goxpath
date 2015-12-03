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

func TestFollowing(t *testing.T) {
	p := `//p3/following::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3/><p4/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p4></p4>`, `<p5><p6></p6></p5>`, `<p6></p6>`}
	exec(p, x, exp, t)
}

func TestPreceding(t *testing.T) {
	p := `//p6/preceding::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3/><p4/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p2><p3></p3><p4></p4></p2>`, `<p3></p3>`, `<p4></p4>`}
	exec(p, x, exp, t)
}

func TestPrecedingSibling(t *testing.T) {
	p := `//p4/preceding-sibling::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3><p31/></p3><p4/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p3><p31></p31></p3>`}
	exec(p, x, exp, t)
}

func TestFollowingSibling(t *testing.T) {
	p := `//p2/following-sibling::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3/><p4/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p5><p6></p6></p5>`}
	exec(p, x, exp, t)
}

func TestComment(t *testing.T) {
	p := `//comment()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><!-- comment --></p1>`
	exp := []string{`<!-- comment -->`}
	exec(p, x, exp, t)
}

func TestText(t *testing.T) {
	p := `//text()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	exp := []string{`text`}
	exec(p, x, exp, t)
}

func TestProcInst(t *testing.T) {
	p := `//processing-instruction()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc?></p1>`
	exp := []string{`<?proc?>`}
	exec(p, x, exp, t)
}

func TestProcInst2(t *testing.T) {
	p := `//processing-instruction('proc2')`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc1?><?proc2?></p1>`
	exp := []string{`<?proc2?>`}
	exec(p, x, exp, t)
}

func TestNamespace(t *testing.T) {
	p := `//namespace::*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{`<?namespace http://foo.bar?>`, `<?namespace http://test?>`}
	exec(p, x, exp, t)
}

func TestNamespace2(t *testing.T) {
	p := `//namespace::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:foo="http://test"><p2 xmlns:test="http://foo.bar"></p2></p1>`
	exp := []string{`<?namespace http://foo.bar?>`}
	exec(p, x, exp, t)
}
