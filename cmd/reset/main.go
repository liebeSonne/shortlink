package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Printf("Start generate struct reset methods\n")
	err := run()
	if err != nil {
		log.Fatalf("error on run generate reset: %v", err)
	}
	fmt.Printf("Complete generate struct reset methods\n")
}
