package boollit

//BoolLit is a boolean XPath result
type BoolLit bool

func (b BoolLit) ResValue() string {
	if b {
		return "true"
	}

	return "false"
}
