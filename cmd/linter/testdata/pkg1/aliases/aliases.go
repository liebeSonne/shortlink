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

func (c *customLog) Fatalf(format string, v ...any) {
	fmt.Printf("customLog Fatalf "+format, v)
}

func (c *customLog) Fatalln(v ...any) {
	fmt.Print("customLog Fatalln %v\n", v)
}

func main() {
	os := &customOS{}
	os.Exit(0)
	log := &customLog{}
	log.Fatal("ok")
	log.Fatalf("%d", 0)
	log.Fatalln("ok")
	aliaslog.Fatal("ok")
	aliaslog.Fatalf("ok %d", 1)
	aliaslog.Fatalln("ok")
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
	log.Fatalf("%d", 0)
	aliaslog.Fatal("warning")        // want "log.Fatal call in main package not in main function is not allowed"
	aliaslog.Fatalf("warning %d", 1) // want "log.Fatalf call in main package not in main function is not allowed"
	aliaslog.Fatalln("warning")      // want "log.Fatalln call in main package not in main function is not allowed"
}

func DoSomethingWithOsExit() {
	os := &customOS{}
	os.Exit(0)
	aliasos.Exit(0) // want "os.Exit call in main package not in main function is not allowed"
}
