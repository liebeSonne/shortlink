package complex

// generate:reset
type ComplexStruct struct {
	BoolField    bool
	IntField     int
	StringField  string
	FloatField   float64
	ByteField    byte
	RuneField    rune
	ComplexField complex128
	SliceField   []int
	MapField     map[string]int
	ArrayField   [3]int
	ChanField    chan int
	FuncField    func()
	StructPtr    *SubStruct
}

type SubStruct struct {
	Value int
}
