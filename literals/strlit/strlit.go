package strlit

//StrLit is a string XPath result
type StrLit string

//ResValue implements the tree.Res interface so strings can be returned from XPath expressions.
func (s StrLit) ResValue() string {
	return string(s)
}
