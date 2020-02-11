package compat

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/dwood15/bosskeys/bk"
)

type OotRregion struct {
	RegionName string
	Locations  map[string]string
	Exits      map[string]string
}

func (otr OotRregion) ToNodeChunk() (nl []bk.Node) {
	//A region is a node of class hub
	var rNode bk.Node
	rNode.Name = bk.NodeName(otr.RegionName)
	rNode.Class = bk.Hub

	numNew := len(otr.Locations) + len(otr.Exits)
	rNode.Exits = make([]string, numNew)
	nl = make([]bk.Node, numNew)

	//A Location is a loopback - it's only connected to the parent area, and once visited,
	//with requirements _met_, self-destructs, so it can't be visited by the search algo
	//again, helping improve search performance as it runs.
	// ALL locations give something
	for k, v := range otr.Locations {
		rNode.Exits = append(rNode.Exits, bk.NodeName(k))

		n := bk.Node{
			Name:     bk.NodeName(k),
			Class:    bk.OneWayPortal,
			Requires: []bk.KeyName{bk.KeyName(v)},
			OnVisit: &struct {
				Gives         []bk.KeyName
				SelfDestructs bool
			}{
				Gives: []bk.KeyName{bk.KeyName("HCKING TODO")},
				SelfDestructs: true,
			},
			Exits: []bk.NodeName{ rNode.Name },
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
				Gives: []bk.KeyName{bk.KeyName("HCKING TODO")},
				SelfDestructs: false,
			},
			Exits: []bk.NodeName{ rNode.Name },
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

func loadRegions(filename string) (regs []*OotRregion) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, &regs); err != nil {
		panic(err)
	}

	if len(regs) == 0 {
		panic("Loaded Regions is zero!!")
	}

	for _, r := range regs {
		r.RegionName = toLowerSnake(r.RegionName)
		r.Locations = mToLowerSnake(r.Locations)
		r.Exits = mToLowerSnake(r.Exits)
	}

	return

}

func ConvertOOT(wd string) {
	regs := loadRegions(wd + "overworld.json")
	ns := make([]bk.Node, len(regs))

	println("loaded :", len(regs), " regions. converting to nodes")

}
