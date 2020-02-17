package bk

import (
	"fmt"
)

type (
	NodeName  string //NodeName is the human-readable name of the node
	NodeClass string //NodeClass represents a category of node

	//helper collections to make searching through them easier
	NodeClasses []NodeClass
)

const (
	Loopback     NodeClass = "loopback"       // Blue Warps and Owl teleport
	OneWayPortal NodeClass = "one_way_portal" // Blue Warps and Owl teleport
	Hub          NodeClass = "hub"            // Hubs may contain items and exits
)

type (
	OnVisit struct {
		Gives         []KeyName `json:"gives,omitempty"`          //Gives is a list of Human-Readable items
		SelfDestructs bool      `json:"self_destructs,omitempty"` //Whether or not this node self-destructs after visiting and taking the associated item
	}

	DistributionSettings struct {
		ItemLocked     bool `json:"items_locked,omitempty"` // When items are locked, this means that, the items here DO NOT enter the shuffle pool
		EntranceLocked bool `json:"omitempty"`              // EntranceLocked - not doing entrance randomizer right now
	}

	NodeList []Node
	Node     struct {
		Name         NodeName             `json:"name,omitempty"` // Name is the human-readable identifier of the particular Node.
		Comment      string               `json:"comments,omitempty"`
		MiniMapScene string               `json:"mini_map_scene,omitempty"`
		Class        NodeClass            `json:"class,omitempty"`    // Class is a descriptor of the node
		Requires     KeyPhrase            `json:"requires,omitempty"` // Names of the Items/Flags that are required in order to visit this node.
		OnVisit      *OnVisit             `json:"on_visit,omitempty"`
		Exits        []NodeName           `json:"exits,omitempty"`
		Settings     DistributionSettings `json:"settings,omitempty"`

		destructed bool
		index      int //index is where it is in the NodeList array//pool
	}
)

//NewNode is a hack for some mild awkwardness with embedded structs-to-pointers
func NewNode() Node {
	return Node{
		OnVisit: &OnVisit{},
	}
}

//Validation helpers
var AllNodeClasses = NodeClasses{OneWayPortal, Loopback, Hub}

//Major helper funcs

//CanVisit indicates whether or not we are able to access the next node and therefore claim a given item
func (n *Node) CanVisit(from *Node, keysHeld map[KeyName]Key) bool {

	//Sanity checks. We panic because Programmer error means the local node is empty or missing things
	if n == nil {
		panic("can't check visit for nil nodes")
	}

	if len(n.Name) == 0 {
		panic("invalid node Name - can't ")
	}

	//keysHeld check assumes that during testing, the algo has at least one key.
	if len(keysHeld) == 0 {
		println("rejecting visit to node: " + n.Name)
		return false
	}

	if len(from.Name) == 0 {
		panic("'from' node should never be emptystring")
	}

	if len(n.Requires) == 0 {
		return true
	}

	//idea: return items which the algoshould think are missing?

	n.Requires.ParseRequirements()
	return false

	//k, ok := keysHeld[n.Requires]
	//if !ok || len(k.Name) == 0 {
	//	return false
	//}
	//
	//if !k.Use(keysHeld) {
	//	return false
	//}
	//
	////golang's funky about modifying members of a map...
	////I'm a scrub so we reassign it back to the map
	//keysHeld[req] = k

	return true
}

func (n *Node) Validate() error {
	if len(n.Name) == 0 {
		return fmt.Errorf("no name. cannot use for tree traversal")
	}

	if !AllNodeClasses.Contains(string(n.Class)) {
		return fmt.Errorf("node class: [%s]", n.Class)
	}

	//TODO: More validation of nodes for sanity checking

	switch n.Class {
	case Loopback:
		if n.OnVisit == nil {
			return fmt.Errorf("[%s] missing on_visit - all loopbacks require an on_visit entry", n.Name)
		}

		if len(n.OnVisit.Gives) != 1 {
			return fmt.Errorf("[%s] doesn't have correct number of Gives for class: [%s]", n.Name, n.Class)
		}

		if len(n.Exits) > 0 {
			return fmt.Errorf("[%s] cannot have exits, and yet it still has some: [%v]", n.Name, n.Exits)
		}
	case Hub:
		if len(n.Exits) == 0 {
			return fmt.Errorf("[%s] doesn't have any Exits for class: [%s]", n.Name, n.Class)
		}
	case OneWayPortal:
		if len(n.Exits) != 1 {
			return fmt.Errorf("[%s] doesn't have correct number of Exits for class of: [%s]", n.Name, n.Class)
		}
	default:
		return fmt.Errorf("[%s] has an unrecognized class: [%s]", n.Name, n.Class)
	}

	if n.OnVisit != nil && n.OnVisit.SelfDestructs {
		if len(n.Exits) > 1 {
			return fmt.Errorf("[%s] has too many exits for a self-destructing node. [%s]", n.Name, n.Exits)
		}
	}


	return nil
}

//Minor helper-funcs

//The major issue with golang: no nice generics. :eye_roll:
func (nc NodeClasses) Contains(n string) bool {
	for _, v := range nc {
		if string(v) == n {
			return true
		}
	}

	return false
}
