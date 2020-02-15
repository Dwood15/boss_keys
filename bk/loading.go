package bk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//fixture files - edit these

//helper struct to assist tests
type (
	Pools struct {
		Items       []*Key
		Flags       []*Key
		Nodes       []*Node
		NodesByName map[NodeName]*Node
	}
)

//LoadBasePools pulls the pools from the containing game's folder. At this time, only oot
// is recognized, however others may be added in the future.
func LoadBasePools(wd string) (kg Pools) {
	//sorry windows users :P
	if len(wd) == 0 {
		kg.Items = LoadKeyPool("bk/base_pools/oot/item_pool.json")
		kg.Flags = LoadKeyPool("bk/base_pools/oot/state_flags.json")
	}
	//Just one ginormous json file -- rip.
	kg.Nodes = LoadNodes(wd + "nodes.json")
	kg.NodesByName = make(map[NodeName]*Node, len(kg.Nodes))
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

func LoadAndValidateNodes(wd string) (errs []error) {
	kg := LoadBasePools(wd)
	nl := kg.Nodes

	println("validating bk nodes. Num: ", len(nl))
	errs = make([]error, 0, len(nl))

	for i, n := range nl {
		var err error
		//Check that the internal data is coherent
		if err = n.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("Node Idx: [%d], err: [%s]", i, err.Error()))
			continue
		}
		//Check that there is only one region of any given name
		_, ok := kg.NodesByName[n.Name]
		if ok {
			errs = append(errs, fmt.Errorf("Node Idx: [%d], err: [%s] already exists", n.Name))
		}
	}

	return errs
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
