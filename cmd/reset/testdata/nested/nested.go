package nested

// generate:reset
type Parent struct {
	Child      Child
	ChildPtr   *Child
	GrandChild *GrandChild
}

// generate:reset
type Child struct {
	Value string
}

// generate:reset
type GrandChild struct {
	Data int
}
