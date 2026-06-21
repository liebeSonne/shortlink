package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Printf("Start generate struct reset methods\n")
	err := run()
	if err != nil {
		log.Fatal(fmt.Errorf("error on run generate reset: %w", err))
	}
	fmt.Printf("Complete generate struct reset methods\n")
}
