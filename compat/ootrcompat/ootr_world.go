package compat

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/dwood15/bosskeys/bk"
)

type OotRregion struct {
	RegionName string `json:"region_name"`
	Locations  map[string]string
	Exits      map[string]string
}

type OotItem map[string]bk.KeyName

func (otr OotRregion) ToNodeChunk(itemsByLoc *map[string]string) (nl []bk.Node) {
	//A region is a node of class hub
	var rNode bk.Node
	rNode.Name = bk.NodeName(otr.RegionName)
	rNode.Class = bk.Hub

	numNew := len(otr.Locations) + len(otr.Exits)
	rNode.Exits = make([]bk.NodeName, 0, numNew)
	nl = make([]bk.Node, 0, numNew)

	//A Location is a loopback - it's only connected to the parent area, and once visited,
	//with requirements _met_, self-destructs, so it can't be visited by the search algo
	//again, helping improve search performance as it runs.
	// ALL locations give something
	for k, v := range otr.Locations {
		rNode.Exits = append(rNode.Exits, bk.NodeName(k))

		var n bk.Node

		n = bk.Node{
			Name:     bk.NodeName(k),
			Class:    bk.OneWayPortal,
			Comment:  "",
			Requires: []bk.KeyName{bk.KeyName(v)},
			Exits:    []bk.NodeName{rNode.Name},
		}

		itmName, ok := (*itemsByLoc)[k]
		if ok {
			n.OnVisit.Gives = []bk.KeyName{bk.KeyName(itmName)}
			n.OnVisit.SelfDestructs = true
			delete(*itemsByLoc, itmName)
		}

		nl = append(nl, n)
	}

	//An Exit is considered the same as a Location, a repeatable portal
	for k, v := range otr.Exits {
		rNode.Exits = append(rNode.Exits, bk.NodeName(k))

		n := bk.Node{
			Name:     bk.NodeName(k),
			Class:    bk.OneWayPortal,
			Requires: []bk.KeyName{bk.KeyName(v)},
			OnVisit: &struct {
				Gives         []bk.KeyName
				SelfDestructs bool
			}{
				Gives:         []bk.KeyName{bk.KeyName("FRCKING TODO")},
				SelfDestructs: false,
			},
			Exits: []bk.NodeName{rNode.Name},
		}

		nl = append(nl, n)
	}

	return nl
}

//spaces become underscores. dumb conversion
func toLowerSnake(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}

func mToLowerSnake(m map[string]string) map[string]string {
	newM := map[string]string{}
	for k, v := range m {
		k = toLowerSnake(k)
		newM[k] = strings.ToLower(v)
	}

	return newM
}

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

var itemMaps = []string{"vanilla_location_items", "shop", "gs_tokens", "dungeon_items"}

func loadItems(wd string) map[string]string {
	itms := make(map[string]string)
	println("attempting to load all the items")

	for _, fName := range itemMaps {
		var tmp map[string]string

		b, err := ioutil.ReadFile(wd + fName + ".json")
		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(b, &tmp); err != nil {
			panic(err)
		}

		if len(tmp) == 0 {
			panic("Loaded items locations map is empty!!")
		}

		for k, v := range tmp {
			itms[toLowerSnake(k)] = toLowerSnake(v)
		}

	}

	newItms := make(map[string]string)

	for k, v := range itms {
		newItms[toLowerSnake(k)] = toLowerSnake(v)
	}

	return newItms
}

func ConvertOOTR(wd string) {

	itms := loadItems(wd)
	regs := loadRegions(wd)
	println("loaded :", len(regs), " regions. converting to nodes")

	ns := make([]bk.Node, 0, len(regs))

	for _, r := range regs {
		ns = append(ns, r.ToNodeChunk(&itms)...)
	}

	if len(itms) != 0 {
		print("num remaining: ", len(itms))
		for k := range itms {
			println("Item Location: [", k, "] not found in loaded regions")
		}
	}

	println("nodes completed and appended: ", len(ns), "total nodes now exist. dumping to file")

	b, err := json.MarshalIndent(ns, "", "  ")

	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile("tmp_oot_nodes.json", b, 0644); err != nil {
		panic(err)
	}

	println("ootr dumps back to file are complete!")
}
