package main

import (
	"log"
	"os"
)

func MockDoSomethingWithPanic() {
	panic("ok")
}

func MockDoSomethingWithLofFatal() {
	log.Fatal("ok")
	log.Fatalf("ok %d", 1)
	log.Fatalln("ok")
}

func MockDoSomethingWithOsExit() {
	os.Exit(0)
}
