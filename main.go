package main

import (
	"io/ioutil"

	"github.com/dwood15/bosskeys/bk"

	"encoding/json"

	"github.com/alecthomas/jsonschema"
)

func main() {
	println("Launching interactive terminal for building json")
	println("for starters, loading all of the pools.")

	bk.LoadBasePools("bk/base_pools/oot/")

	println("Dumping Json schema to file")



	s := jsonschema.Reflect(&[]bk.Node{})

	b, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("node_schema.jsonschema", b, 0644)
	if err != nil {
		panic(err)
	}

	println("node_schema json schema output")
}
