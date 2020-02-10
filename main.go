package main

import (
	"fmt"

	"github.com/dwood15/bosskeys/bk"
)

func main() {
	println("Launching interactive terminal for building json")
	println("for starters, loading all of the pools.")

	bk.LoadBasePools("bk/base_pools/oot/")

	println("Pools loaded. Edit deku_tree.json")

	var input string
	var exit bool
	for !exit {
		println("waiting for input:")
		_, err := fmt.Scanln(&input)

		if err != nil {
			println("wtf, this shouldn't happen. err: ", err.Error())
			return
		}

		println("you entered: ", input)
		if input == "x" {
			println("exiting")
			return
		}
	}
}
