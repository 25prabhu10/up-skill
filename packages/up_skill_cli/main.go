// Package main is the entry point of the up-skill CLI application. It initializes and executes the root command.
package main

import (
	"os"

	"github.com/25prabhu10/up-skill/cmd/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
