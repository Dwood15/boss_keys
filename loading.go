package bosskeys

import (
	"encoding/json"
	"io/ioutil"
)

//fixture files - edit these

//helper struct to assist tests
type Pools struct {
	Items []*Key
	Flags []*Key

	Nodes []*Node
}

//LoadBasePools pulls the pools from the containing game's folder. At this time, only oot
// is recognized, however others may be added in the future.
func LoadBasePools(wd string) (kg Pools) {
	//sorry windows users :P
	kg.Items = LoadKeyPool(wd + "item_pool.json")
	kg.Flags = LoadKeyPool(wd + "state_flags.json")
	kg.Nodes = LoadNodes(wd + "nodes.json")

	return
}

func LoadKeyPool(filename string) (keys []*Key) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, &keys); err != nil {
		panic(err)
	}

	if len(keys) == 0 {
		panic("keys list is zero!!")
	}

	for _, k := range keys {
		if err = k.Validate(); err != nil {
			panic(err)
		}
	}

	return
}

func LoadNodes(filename string) (nl []*Node) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, &nl); err != nil {
		panic(err)
	}

	if len(nl) == 0 {
		panic("keys list is zero!!")
	}

	for _, n := range nl {
		if err = n.Validate(); err != nil {
			panic(err)
		}
	}

	return
}
