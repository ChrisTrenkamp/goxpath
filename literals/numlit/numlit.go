package numlit

import "fmt"

//NumLit is a numeric XPath result
type NumLit float64

//ResValue implements the tree.Res interface so numbers can be returned from XPath expressions.
func (n NumLit) ResValue() string {
	return fmt.Sprintf("%g", float64(n))
}
