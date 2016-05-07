package strlit

//StrLit is a string XPath result
type StrLit string

func (s StrLit) ResValue() string {
	return string(s)
}
