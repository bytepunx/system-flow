// Command flai manages monorepos that follow the system-flow standard.
package main

import (
	"os"

	"github.com/bytepunx/system-flow/flai/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
