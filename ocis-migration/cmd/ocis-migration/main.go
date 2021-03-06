package main

import (
	"os"

	"github.com/refs/ocis-mono/ocis-migration/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
