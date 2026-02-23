package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/liebeSonne/shortlink/internal/config"
)

const defaultURLAddress = "http://localhost:8080"
const defaultServerAddress = ":8080"

var ErrInvalidFlagValue = errors.New("invalid flag value")
var ErrInvalidDefaultServerAddress = errors.New("invalid default server address")

func parseFlags(config *config.Config) error {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)

	serverAddress := address{}
	err := serverAddress.Set(defaultServerAddress)
	if err != nil {
		log.Printf("invalid default server address: %v", err)
		return ErrInvalidDefaultServerAddress
	}

	fs.Var(&serverAddress, "a", "address and port to run server")
	urlAddress := fs.String("b", defaultURLAddress, "address and port for output short url")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		log.Printf("error parsing config flags: %v", err)
		return err
	}

	config.ServerAddress = serverAddress.String()
	config.URLAddress = *urlAddress

	return nil
}

type address struct {
	Host string
	Port int
}

func (a *address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
func (a *address) Set(flagValue string) error {
	params := strings.Split(flagValue, ":")

	if len(params) != 2 {
		return ErrInvalidFlagValue
	}

	port, err := strconv.Atoi(params[1])
	if err != nil {
		log.Printf("error on atoi port: %v\n", err)
		return ErrInvalidFlagValue
	}

	a.Host = params[0]
	a.Port = port
	return nil
}
