package main

import (
	"github.com/cjs/purtypics/cmd"
)

// version is set by build flags
var version = "dev"

func main() {
	cmd.Execute(version)
}