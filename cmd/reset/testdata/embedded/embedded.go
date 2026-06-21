package embedded

// generate:reset
type Embedded struct {
	Value int
}

// generate:reset
type WithEmbedded struct {
	Embedded // встроенное поле
	Name     string
}
