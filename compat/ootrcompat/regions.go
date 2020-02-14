package ootrcompat

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/dwood15/bosskeys/bk"
)

var regFiles = []string{"overworld", "deku_tree", "botw", "dd_cavern", "fire", "forest", "ganon", "ice_cavern", "jj_belly", "shadow", "spirit", "training_grounds", "water"}

const ext = ".json"

func loadRegions(wd string) (regs []*OotRregion) {

	for _, v := range regFiles {
		b, err := ioutil.ReadFile(wd + v + ext)
		if err != nil {
			panic(err)
		}

		var rL []*OotRregion
		if err = json.Unmarshal(b, &rL); err != nil {
			println("error on file: ", v)
			panic(err)
		}

		if len(rL) == 0 {
			panic("Loaded Regions is zero!!")
		}

		regs = append(regs, rL...)
	}

	//Pre-Sanitizing
	for _, r := range regs {
		r.RegionName = toLowerSnake(r.RegionName)
		r.Locations = mToLowerSnake(r.Locations)
		r.Exits = mToLowerSnake(r.Exits)
	}

	return
}

func mToLowerSnake(m OotRLocations) OotRLocations {
	newM := OotRLocations{}
	for k, v := range m {
		k = bk.NodeName(toLowerSnake(string(k)))
		newM[k] = OotRRequirement(strings.ToLower(string(v)))
	}

	return newM
}