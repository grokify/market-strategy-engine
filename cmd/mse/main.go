// Command mse is the CLI for the Market Strategy Engine.
//
// Usage:
//
//	mse [command]
//
// Available Commands:
//
//	generate-schema  Generate JSON Schema files from Go types
//	validate         Validate an analysis file against schema
//	help             Help about any command
package main

import (
	"os"

	"github.com/grokify/market-strategy-engine/cmd/mse/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
