package pkg1

import (
	"fmt"
	"log"
	"os"
)

func DoSomething() {
	fmt.Println("ok")
}

func DoSomethingWithPanic() {
	panic("warning") // want "panic call is not allowed"
}

func DoSomethingWithLofFatal() {
	log.Fatal("ok")
	log.Fatalf("ok %d", 1)
	log.Fatalln("ok")
}

func DoSomethingWithOsExit() {
	os.Exit(0)
}
