package strlit

//StrLit is a string XPath result
type StrLit string

func (s StrLit) String() string {
	return string(s)
}
