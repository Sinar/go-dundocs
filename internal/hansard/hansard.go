package hansard

type HansardType int

const (
	LISAN HansardType = iota
	BUKANLISAN
	DEBAT
)

type Hansard struct{}
