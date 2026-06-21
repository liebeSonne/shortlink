package testdata

// generate:reset
type User struct {
	ID       int
	Name     string
	Email    string
	Tags     []string
	Metadata map[string]interface{}
	Profile  *Profile
	Child    *Child
}

// generate:reset
type Profile struct {
	Age  int
	City string
	Bio  string
}

// generate:reset
type Product struct {
	SKU   string
	Price float64
	Stock int
}

// generate:reset
type Child struct {
	Value string
}

type NoResetStruct struct {
	Value int
}
