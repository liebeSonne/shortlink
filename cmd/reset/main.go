package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.Info("Start generate struct reset methods")
	err := run()
	if err != nil {
		slog.Error("error on run generate reset", "error", err)
		os.Exit(1)
	}
	slog.Info("Complete generate struct reset methods")
}
