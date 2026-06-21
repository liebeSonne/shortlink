package collections

// generate:reset
type CollectionStruct struct {
	IntSlice       []int
	StringSlice    []string
	StringMap      map[string]string
	IntMap         map[int]int
	StructSlice    []SubStruct
	StructMap      map[string]SubStruct
	InterfaceSlice []interface{}
	InterfaceMap   map[string]interface{}
}

type SubStruct struct {
	ID   int
	Name string
}
