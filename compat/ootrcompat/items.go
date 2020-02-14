package ootrcompat

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dwood15/bosskeys/bk"
)

var itemMaps = []string{"vanilla_location_items", "shop", "gs_tokens",  "scrubs", "dungeon", "event_based", "trade_quest", "drops"}

func loadItems(wd string) OoTRItems {
	itms := make(OotRLocations)
	println("attempting to load all the items")

	for _, fName := range itemMaps {
		var tmp map[string]string

		b, err := ioutil.ReadFile(wd + fName + ".json")
		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(b, &tmp); err != nil {
			panic(fName + " " + err.Error())
		}

		if len(tmp) == 0 {
			panic("Loaded items locations map is empty!!")
		}

		for k, v := range tmp {
			itms[bk.NodeName(toLowerSnake(k))] = OotRRequirement(toLowerSnake(v))
		}

	}

	newItms := make(OoTRItems)

	for k, v := range itms {
		newItms[bk.NodeName(toLowerSnake(string(k)))] = bk.KeyName(toLowerSnake(string(v)))
	}

	return newItms
}