package bk

import (
	"fmt"
)

//helper struct to assist tests
type (
	State struct {
		Inventory map[KeyName]*Key
	}

	Pools struct {
		Items       []*Key
		Flags       []*Key
		Nodes       []*Node
		NodesByName map[NodeName]*Node
	}
)


func (p *Pools) Lookup(nN NodeName) (lUp *Node) {
	return p.NodesByName[nN]
}

//GiveToPlayer plops the item from the Node n and evaluates the on_give s-expr.
//At this time, it does not give anything to the player. TODO: Give the item to the player
func (p *Pools) GiveToPlayer(n *Node) {
	if n.OnVisit != nil {
		if len(n.OnVisit.Gives) > 0 {
			return
		}
	}

	return
}

func (p *Pools) HandleVisit(from *Node, exitIdx int) (to *Node) {
	if len(from.Exits) == 0 {
		return from
	}

	if to = p.Lookup(from.Exits[exitIdx]); to == nil {
		return from
	}

	//TODO: After a node is destroyed, ensure it can be mapped.
	if to.destructed {
		return from
	}


	if !to.CanVisit(from, nil) {
		return from
	}

	p.GiveToPlayer(to)

	if numExits := len(to.Exits); numExits == 0 {
		//No validation of class + exit pair during computation time.
		return from
	}

	return to
}

//Destruct handles the self-destruct which happens after an OnVisit. it will zip up any exits to this
func (p *Pools) Destruct(parent *Node, n*Node) {
	//TODO: actually "destroy* the current node. But first, it must know what hooks to it.
	//For now, we will simply set the "Destroyed" flag
	n.destructed = true

	// pluck the node from the slice.
	//golang has no pop() mechanism, *sigh*
	//p.Nodes[n.index] = nil

	// remove the node from the lookupcache.
	//delete(p.NodesByName, n.Name)
}

//FastTraverse travels the tree from a given start location to a target.
//It skips Flags and Items and OnGive events. The name may be a misnomer.
//It is only optimized for complexity, not runtime
func (p *Pools) FastTraverse(start, target NodeName) {
	sNode := p.NodesByName[start]
	tNode := p.NodesByName[target]

	if sNode == nil || tNode == nil {
		panic(fmt.Sprintf("invalid/missing node name(s) specified: either [%s] or [%s]", start, target))
	}

	currentNode := sNode

	for exitIdx, e := range currentNode.Exits {
		eN := p.NodesByName[e]

		p.HandleVisit(eN, exitIdx)

		//Nodes that self-destruct _should_ be loopback nodes, however that may change.
		if eN.OnVisit != nil && eN.OnVisit.SelfDestructs {
			p.Destruct(currentNode, eN)
		}
	}
}