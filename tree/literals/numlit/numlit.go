package numlit

import "fmt"

//NumLit is a numeric XPath result
type NumLit float64

func (n NumLit) String() string {
	return fmt.Sprintf("%g", float64(n))
}
