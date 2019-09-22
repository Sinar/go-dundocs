package hansard

type HansardType int

const (
	HANSARD_SPOKEN HansardType = iota
	HANSARD_WRITTEN
	HANSARD_DEBATE
)

type Hansard struct{}
