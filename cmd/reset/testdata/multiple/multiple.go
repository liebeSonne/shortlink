package multiple

// generate:reset
type User struct {
	ID    int
	Name  string
	Age   int
	Email string
}

// generate:reset
type Admin struct {
	User
	Role        string
	Permissions []string
}

// generate:reset
type Guest struct {
	ID       int
	Settings map[string]interface{}
}
