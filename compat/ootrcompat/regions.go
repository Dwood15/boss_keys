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

	println("regions loaded, processing to lower snake and calculating cow locations")
	var cows = make(map[bk.NodeName]string)

	//Pre-Sanitizing
	for _, r := range regs {
		for k := range r.Locations {
			if string(k) == "Impas House Near Cow" {
				continue
			}

			if strings.Contains(string(k), "Cow") {
				cows[k] = "Milk"
			}
		}

		r.RegionName = toLowerSnake(r.RegionName)
		r.Locations = mToLowerSnake(r.Locations)
		r.Exits = mToLowerSnake(r.Exits)
	}

	cowBytes, _ := json.MarshalIndent(cows, "", "  ")
	_ = ioutil.WriteFile("cows.json", cowBytes, 0644)

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