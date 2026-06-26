package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.Fatal("ok")
	log.Fatalf("ok: %d", 1)
	log.Fatalln("ok")
	os.Exit(0)
	panic("warning") // want "panic call is not allowed"
}

func DoSomething() {
	fmt.Println("ok")
}

func DoSomethingWithPanic() {
	panic("warning") // want "panic call is not allowed"
}

func DoSomethingWithLofFatal() {
	log.Fatal("warning")        // want "log.Fatal call in main package not in main function is not allowed"
	log.Fatalf("warning %d", 1) // want "log.Fatalf call in main package not in main function is not allowed"
	log.Fatalln("warning")      // want "log.Fatalln call in main package not in main function is not allowed"
}

func DoSomethingWithOsExit() {
	os.Exit(0) // want "os.Exit call in main package not in main function is not allowed"
}
