package boollit

//BoolLit is a boolean XPath result
type BoolLit bool

//ResValue implements the tree.Res interface so boolean's can be returned from XPath expressions.
func (b BoolLit) ResValue() string {
	if b {
		return "true"
	}

	return "false"
}
