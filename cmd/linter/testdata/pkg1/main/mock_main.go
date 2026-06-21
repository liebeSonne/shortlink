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
}

func MockDoSomethingWithOsExit() {
	os.Exit(0)
}
