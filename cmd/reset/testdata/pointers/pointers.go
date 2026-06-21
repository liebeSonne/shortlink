package pointers

// generate:reset
type PointerStruct struct {
	IntPtr    *int
	StringPtr *string
	BoolPtr   *bool
	StructPtr *SubStruct
}

// generate:reset
type SubStruct struct {
	Value int
	Name  string
}
