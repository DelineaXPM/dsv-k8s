//go:build ignore

// This file lets you run mage with a no-install option as long as you have go.
// To invoke just run go run main.go [task] [parameters]
// To use mage directly, install it, then run mage [task] [parameters]
package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() { os.Exit(mage.Main()) }
