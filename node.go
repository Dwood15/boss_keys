package bosskeys

import (
	"fmt"
)

type (
	NodeName  string //NodeName is the human-readable name of the node
	NodeClass string //NodeClass represents a category of node

	KeyName      string
	KeyCondition string //KeyCondition represents a requirement for using an item. A KeyCondition is either can_act, or the name of another key
	KeyAction    string //KeyAction indicates what to do after use of the key

	Action string //Action represents what to do when this node is visited

	//helper collections to make searching through them easier
	NodeClasses []NodeClass
	Actions     []Action
	KeyActions  []KeyAction
)

const (
	OneWayPortal NodeClass = "one_way_portal" // Blue Warps and Owl teleport
	TwoWayPortal NodeClass = "two_way_portal" // Doors, keyed entrances
	SingleGive   NodeClass = "single_give"    // Chests, GS, freestanding items
	Toggle       NodeClass = "toggle"         // Child -> Adult, visa versa
	Hub          NodeClass = "hub"            // Hubs may contain items and exits
	Special      NodeClass = "special"
	Interior     NodeClass = "interior" // An interior has one exit. May contain multiple items

	OnUseDoNothing KeyAction = "do_nothing"
	OnUseDecrement KeyAction = "decrement"
	OnUseTeleport  KeyAction = "teleport_to"
)

type (
	Node struct {
		Name     NodeName  // Name is the human-readable identifier of the particular Node.
		Class    NodeClass // Class is a descriptor of the node
		Requires []KeyName // Names of the Items/Flags that are required in order to visit this node.
		OnVisit  *struct {
			Gives         []KeyName //Gives is a list of Human-Readable items
			SelfDestructs bool      //Whether or not this node self-destructs after visiting and taking the associated item
		}

		Exits []string
	}

	// Key represents game state, or player save file state. Anything that can be used to indicate progression, really.
	Key struct {
		Name       KeyName        // Name is the human-readable ID of this key.
		Type       string         // Type is an extra descriptor for a key that can be added in lieu of listing all required items at once
		Conditions []KeyCondition // Conditions is a list of requirements in order to use this item. Expexts a KeyName
		State      struct {
			Action     KeyAction // Action: What to do on use of this key
			TeleportTo NodeName  // TeleportTo: Node to visit. Only valid if Action is teleport
			Value      int       // Value: the current number of this key in inventory
		}
	}
)

//Validation helpers
var AllNodeClasses = NodeClasses{OneWayPortal, TwoWayPortal, SingleGive, Toggle, Hub, Interior, Special}
var AllKeyActions = KeyActions{OnUseDecrement, OnUseDoNothing, OnUseTeleport, ""}

//Major helper funcs

//CanVisit indicates whether or not we are able to access the next node and therefore claim a given item
func (n *Node) CanVisit(from NodeName, keysHeld map[KeyName]Key) bool {
	if len(n.Requires) == 0 {
		return true
	}

	//idea: return items which are missing?

	for _, req := range n.Requires {
		k, ok := keysHeld[req]
		if !ok || len(k.Name) == 0 {
			return false
		}

		if !k.Use(keysHeld) {
			return false
		}

		//golang's funky about modifying members of a map...
		//I'm a scrub so we reassign it back to the map
		keysHeld[req] = k
	}

	return true
}

func (k *Key) Use(otherKeys map[KeyName]Key) (success bool) {
	if len(k.Conditions) == 0 {
		goto act
	}

	for _, condKey := range k.Conditions {
		if condKey == "can_act" {
			continue
		}

		//This bit here assumes that in order to use one key, we just have to have met the other key, _not_ used it.
		otherKey, ok := otherKeys[KeyName(condKey)]
		if !ok || otherKey.Validate() != nil {
			return false
		}
	}

act:
	if len(k.State.Action) == 0 {
		panic("invalid action: empty string")
	}

	if k.State.Action == OnUseDoNothing {
		return true
	}

	if k.State.Action == OnUseDecrement {
		if k.State.Value <= 0 {
			return false
		}

		k.State.Value--
		return true
	}

	//This shouldn't happen, I think?
	return false
}

//Basic sanity checks
func (k *Key) Validate() error {
	if len(k.Name) == 0 {
		return fmt.Errorf("all keys must have a name")
	}

	if !AllKeyActions.Contains(string(k.State.Action)) {
		return fmt.Errorf("key action: [%s] is invalid. must be from predeclared list", k.State.Action)
	}

	if k.State.Action == OnUseTeleport && len(k.State.TeleportTo) == 0 {
		return fmt.Errorf("TeleportTo must be ")
	}

	return nil
}

func (n *Node) Validate() error {
	if len(n.Name) == 0 {
		return fmt.Errorf("node has no name. cannot use for tree traversal")
	}

	if !AllNodeClasses.Contains(string(n.Class)) {
		return fmt.Errorf("node class: [%s]", n.Class)
	}

	//TODO: More validation of nodes for sanity checking

	switch n.Class {
	case SingleGive:
		if len(n.OnVisit.Gives) != 1 {
			return fmt.Errorf("node [%s] doesn't have correct number of Gives for class: [%s]", n.Name, n.Class)
		}
	case Toggle:
		panic("not implemented yet!")
	case TwoWayPortal:
		if len(n.Exits) != 2 {
			return fmt.Errorf("node [%s] doesn't have correct number of Exits for class: [%s]", n.Name, n.Class)
		}
	case OneWayPortal:
		if len(n.Exits) != 1 {
			return fmt.Errorf("node [%s] doesn't have correct number of Exits for class of: [%s]", n.Name, n.Class)
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

func (a Actions) Contains(n string) bool {
	for _, v := range a {
		if string(v) == n {
			return true
		}
	}

	return false
}

func (a KeyActions) Contains(n string) bool {
	for _, v := range a {
		if string(v) == n {
			return true
		}
	}

	return false
}
