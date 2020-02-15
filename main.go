package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/alecthomas/jsonschema"

	"github.com/dwood15/bosskeys/bk"
	"github.com/dwood15/bosskeys/compat/ootrcompat"
)

func reflectSchemas() {
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

func main() {
	println("Launching interactive terminal for building json")

	//Convert from the ootr logic files to the our new ones
	ootrcompat.ConvertOOTR("compat/ootrcompat/")


	println("everything converted, check")
	_, _ = fmt.Scanln()

	//"bk/base_pools/oot/"
	errs := bk.LoadAndValidateNodes("")

	for _, err := range errs {
		println(err.Error())
	}
}
