package main

import (
	"fmt"
	aliaslog "log"
	aliasos "os"
)

type customOS struct {
}

func (c *customOS) Exit(code int) {
	fmt.Print("customOS Exit %d\n", code)
}

type customLog struct {
}

func (c *customLog) Fatal(v ...any) {
	fmt.Print("customLog Fatal %v\n", v)
}

func main() {
	os := &customOS{}
	os.Exit(0)
	log := &customLog{}
	log.Fatal(0)
	aliaslog.Fatal("ok")
	aliasos.Exit(0)
	panic("warning") // want "panic call is not allowed"
}

func DoSomething() {
	fmt.Println("ok")
}

func DoSomethingWithPanic() {
	panic("warning") // want "panic call is not allowed"
}

func DoSomethingWithLofFatal() {
	log := &customLog{}
	log.Fatal(0)
	aliaslog.Fatal("warning") // want "log.Fatal call in main package not in main function is not allowed"
}

func DoSomethingWithOsExit() {
	os := &customOS{}
	os.Exit(0)
	aliasos.Exit(0) // want "os.Exit call in main package not in main function is not allowed"
}
