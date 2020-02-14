package ootrcompat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dwood15/bosskeys/bk"
)

type OoTRItems map[bk.NodeName]bk.KeyName
type OotRRequirement string
type OotRLocations map[bk.NodeName]OotRRequirement
type OotRregion struct {
	RegionName string `json:"region_name"`
	Scene      string `json:"scene,omitempty"`
	Hint       string `json:"hint,omitempty"`
	Locations  OotRLocations
	Exits      OotRLocations
}

func ConvertOOTR(wd string) {
	itms := loadItems(wd + "items/")

	ogItmCt := len(itms)
	println("loaded: [", ogItmCt, "] itms loaded")
	if ogItmCt == 0 {
		panic("no items were loaded from the static files. this is not acceptable")
	}

	regs := loadRegions(wd + "areas/")
	println("loaded: [", len(regs), "] regions. converting to nodes")

	nL := make([]bk.Node, 0, len(regs))
	ns := make(map[bk.NodeName]int, len(regs))

	for _, r := range regs {
		for _, n := range r.ToNodeChunk(itms) {
			var _n bk.Node
			_ni, ok := ns[n.Name]
			if !ok {
				ns[n.Name] = len(nL)
				nL = append(nL, n)
				continue
			}

			_n = nL[_ni]

			//The algorithm for building locations is busted, and I'm
			//to hcking lazy to fix the root cause, so this hack will forever burden me.
			if n.OnVisit == nil && _n.OnVisit != nil {
				n.OnVisit = _n.OnVisit
			} else if n.OnVisit != nil && _n.OnVisit != nil {
				if len(n.OnVisit.Gives) == 0 {
					n.OnVisit.Gives = _n.OnVisit.Gives
				} else if len(_n.OnVisit.Gives) != 0 {
					n.OnVisit.Gives = append(n.OnVisit.Gives, _n.OnVisit.Gives...)
				}
			}

			if len(_n.Exits) > 0 && len(n.Exits) == 0 {
				n.Exits = _n.Exits
				goto checkReq
			}

			if len(_n.Exits) > 0 {
				for _, _nE := range _n.Exits {
					var exists bool

					for _, nE := range n.Exits {
						if nE == _nE {
							exists = true
						}
					}

					//no duplicate exits
					if !exists {
						n.Exits = append(n.Exits, _nE)
					}
				}

			}

		checkReq:
			if len(n.Requires) == 0 {
				n.Requires = _n.Requires
			}

			nL[_ni] = n
		}
	}

	if len(itms) != 0 {
		print("NOT ALL items were successfully removed from the pool. Num Remaining: ", len(itms))
		for k := range itms {
			println("Item Location: [", k, "] not found in loaded regions")
			println("Press enter when ready to continue")
			_, _ = fmt.Scanln()
		}
	}

	println("nodes completed and appended: ", len(nL), "total nodes now exist. dumping to file")

	b, err := json.MarshalIndent(nL, "", "  ")

	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile("tmp_oot_nodes.json", b, 0644); err != nil {
		panic(err)
	}

	println("ootr dumps back to file are complete!")
}

//bad func sig yes I know, leave me alone (. __ .)
func locationToNode(iBL OoTRItems, rnName, k bk.NodeName, req OotRRequirement, c bk.NodeClass) bk.Node {
	n := bk.NewNode()

	n.Name = k

	n.Class = c

	n.Comment = "automatically generated by ootr_world.go"

	if req != "True" {
		//These casts are pretty horrible, but left here deliberately to remind
		//that the Requires
		n.Requires = bk.KeyPhrase(req)
	}

	n.Exits = []bk.NodeName{rnName}

	itmName, ok := iBL[k]
	if ok {
		n.OnVisit.Gives = []bk.KeyName{itmName}
		n.OnVisit.SelfDestructs = c == bk.OneWayPortal
		delete(iBL, k)
	} else {
		n.OnVisit = nil
	}

	return n
}

func (otr *OotRregion) ToNodeChunk(itemsByLoc OoTRItems) (nl []bk.Node) {
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
		rNode.Exits = append(rNode.Exits, k)

		n := locationToNode(itemsByLoc, rNode.Name, k, v, bk.OneWayPortal)
		nl = append(nl, n)
	}

	//An Exit is considered the same as a Location, a repeatable portal
	for k, v := range otr.Exits {
		rNode.Exits = append(rNode.Exits, k)

		n := locationToNode(itemsByLoc, rNode.Name, k, v, bk.Hub)
		nl = append(nl, n)
	}

	return append(nl, rNode)
}

//spaces become underscores. dumb conversion
func toLowerSnake(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
