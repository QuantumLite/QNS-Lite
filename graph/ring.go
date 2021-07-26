package graph

import (
	"fmt"
	"math"

	"example.com/config"
)

type Ring struct {
	ringSize int
	Nodes    []*Node
	Links    []*Link
	Type     string
}

func (ring *Ring) Build() {
	ring.ringSize = config.GetConfig().GetSize()
	ring.Type = "ring"
	ring.generateNodes()
	//fmt.Println("Node IDs", ring.nodes)
	ring.generateLinks()
}

func (ring *Ring) Clear() {
	links := ring.GetLinks()
	for _, link := range links {
		link.IsActive = false
		link.IsReserved = false
		link.Reservation = -1
	}
}

func (ring *Ring) generateNodes() {
	id := make([]int, 1)
	ring.Nodes = make([]*Node, ring.ringSize)
	//for i := 0; i < ring.ringSize; i++ {
	//	ring.Nodes[i] = make([]*Node, ring.ringSize)
	//}
	var node *Node
	for i := 0; i < ring.ringSize; i++ {
		id[0] = i
		//fmt.Println(id)
		ring.Nodes[i] = new(Node)
		//////// IMPORTANT: CHECK THIS!
		ring.Nodes[i].ID = make([]int, 1)
		node = MakeNode(id, config.GetConfig().GetMemory())
		copy(ring.Nodes[i].ID, node.ID)
		ring.Nodes[i].Memory = node.Memory
		//fmt.Println(MakeNode(id, config.GetConfig().GetMemory()))
		//fmt.Println("Node ID", ring.nodes[i].ID)
	}
}

func (ring *Ring) generateLinks() {
	id := make([]int, 1)
	//lifetime := config.GetConfig().GetLifetime()
	ring.Links = make([]*Link, ring.ringSize)
	var link *Link
	for i := 0; i < ring.ringSize; i++ {
		ring.Links[i] = new(Link)
		ring.Links[i].ID = make([]int, 1)
		id[0] = i
		link = MakeLink(id, 0, false)
		copy(ring.Links[i].ID, link.ID)
		ring.Links[i].Age = link.Age
		ring.Links[i].IsActive = link.IsActive
		ring.Links[i].IsPruned = link.IsPruned
		ring.Links[i].Reservation = link.Reservation
		ring.Links[i].IsReserved = link.IsReserved
		//ring.links[i] = MakeLink(id, lifetime, false)
	}
}

func (ring *Ring) GetLinks() []*Link {
	//links := make([]*Link, ring.ringSize)
	//for i := 0; i < ring.ringSize; i++ {
	//	links[i] = grid.Links[i]
	//}
	return ring.Links
}

func (ring *Ring) GetSize() int {
	return ring.ringSize
}

func (ring *Ring) GetNodeIDs() [][]int {
	ids := make([][]int, ring.ringSize)
	for i := 0; i < ring.ringSize; i++ {
		ids[i] = make([]int, 1)
	}
	//fmt.Println(ids[0])
	for i := 0; i < ring.ringSize; i++ {
		ids[i] = []int{ring.Nodes[i].ID[0]}
	}
	return ids
}

func (ring *Ring) GetNeighbors(node *Node) ([]*Node, bool) {
	neighbors := make([]*Node, 1)
	neighbors[0] = new(Node)
	//neighbors[1] = new(Node)
	//fmt.Println(neighbors)
	for _, n := range neighbors {
		n.ID = make([]int, 1)
	}
	beginning := []int{0}
	end := []int{ring.ringSize - 1}
	var previous, next int
	isNil := true
	if IsEqual(node.ID, beginning) {
		previous = ring.ringSize - 1
		next = node.ID[0] + 1
	} else if IsEqual(node.ID, end) {
		previous = node.ID[0] - 1
		next = 0
	} else {
		previous = node.ID[0] - 1
		next = node.ID[0] + 1
	}
	if ring.GetLinkBetween(node, ring.Nodes[previous]).IsPruned == false {
		isNil = false
		if neighbors[0].Memory == 0 {
			neighbors[0] = ring.Nodes[previous]
		} else {
			neighbors = append(neighbors, ring.Nodes[previous])
		}
	}
	if ring.GetLinkBetween(node, ring.Nodes[next]).IsPruned == false {
		isNil = false
		if neighbors[0].Memory == 0 {
			neighbors[0] = ring.Nodes[next]
		} else {
			neighbors = append(neighbors, ring.Nodes[next])
		}
	}
	return neighbors, isNil
}

func (ring *Ring) GetLinkBetween(n1, n2 *Node) *Link {
	id1 := n1.ID
	id2 := n2.ID
	var x int
	if id1[0] == id2[0] {
		fmt.Println("Inside GetLinkBetween: identical nodes.")
		return nil
	}
	if id1[0] < id2[0] {
		x = id1[0]
	} else {
		x = id2[0]
	}
	//fmt.Println("Inside get links", x)
	return ring.Links[x]
}

func (ring *Ring) Distance(n1, n2 *Node, measure string) int {
	if measure == HOP {
		if n1.ID[0] <= n2.ID[0] {
			return int(math.Min(float64(n2.ID[0]-n1.ID[0]), float64(n1.ID[0]+ring.ringSize-n2.ID[0])))
		} else {
			return int(math.Min(float64(n1.ID[0]-n2.ID[0]), float64(n2.ID[0]+ring.ringSize-n1.ID[0])))
		}
	}
	fmt.Println("Ring - Distance: The input measure type is not known!")
	return -1
}

func (ring *Ring) GetType() string {
	return ring.Type
}
