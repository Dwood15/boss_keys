package main

import (
	"github.com/dwood15/bosskeys/bk"
	"github.com/alecthomas/jsonschema"
)

func main() {
	println("Launching interactive terminal for building json")
	println("for starters, loading all of the pools.")

	bk.LoadBasePools("bk/base_pools/oot/")

	println("Dumping Json schema to file")

	jsonschema.Reflect(&TestUser{})

}
