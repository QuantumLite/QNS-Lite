package path

import (
	"math/rand"
	"time"

	"example.com/graph"
)

type nonObliviousLocal struct {
	network    graph.Topology
	path       Path
	isFinished bool
	curr       *graph.Node
	//src        *graph.Node
	//dest       *graph.Node
}

func (nol *nonObliviousLocal) Build(network graph.Topology) {
	nol.path = make([]*graph.Node, 1)
	nol.path[0] = new(graph.Node)
	nol.curr = new(graph.Node)
	nol.isFinished = false
	nol.network = network
	rand.Seed(time.Now().UTC().UnixNano())
	//mg.path[0] = mg.src
	//mg.curr = mg.src
}

func (nol *nonObliviousLocal) Clear() {
	nol.path = nil
	nol.isFinished = false
	nol.path = make([]*graph.Node, 1)
	nol.path[0] = new(graph.Node)
	//nol.curr = nil
}

func (nol *nonObliviousLocal) Find(src, dest *graph.Node) (Path, []int, []int, bool) {
	nol.curr = src
	nol.add(src)
	cntr := 0
	//mapping := make(map[*graph.Node]int)
	mapping := make([]int, 1)
	isMappingNil := true
	mappingCursor := 0
	options := make([]int, 1)
	counter := 0
	anotherCounter := 0
	for !nol.curr.IsEqual(dest) {
		if cntr >= nol.network.GetSize()*nol.network.GetSize() {
			//fmt.Println("Inside Find. Counter overflow.")
			return nil, nil, nil, true
		}
		cntr = cntr + 1
		next := nol.path[len(nol.path)-1]
		var choices []*graph.Node
		check := true
		for check {
			anotherCounter++
			if anotherCounter >= nol.network.GetSize()*nol.network.GetSize() {
				//fmt.Println("Inside Find. anotherCounter overflow.")
				return nil, nil, nil, true
			}
			if len(nol.path) < 2 {
				temp, tempChoices, _ := nol.next(dest)
				next = temp
				choices = tempChoices
				break
			}
			temp, tempChoices, _ := nol.next(dest)
			next = temp
			choices = tempChoices
			if !graph.IsEqual(next.ID, nol.path[len(nol.path)-2].ID) {
				////////////////////////////////////// Maybe I can do some pruning here.
				check = false
			}
		}
		//fmt.Println("next is", next.ID, "last is", mg.path[len(mg.path)-1])
		if len(choices) > 1 {
			isMappingNil = false
			if mappingCursor == 0 {
				mapping[mappingCursor] = len(nol.path) - 1
			} else {
				mapping = append(mapping, len(nol.path)-1)
			}
			if counter == 0 {
				options[counter] = len(choices)
			} else {
				options = append(options, len(choices))
			}
			//options[counter] = len(choices)
			counter++
			mappingCursor++
		}
		//fmt.Println("Next found", next, "CNTR", cntr)
		//fmt.Println("This is the path after finding:")
		//for _, nodede := range mg.path {
		//	fmt.Println(nodede.ID)
		//}
		if next == nil {
			//fmt.Println("nil next", mg.path)
			return nil, nil, nil, true
		}
		//mg.add(mg.next(dest))
		nol.add(next)
		nol.curr = nol.path[len(nol.path)-1]
	}
	//fmt.Println("Found Path - inside find: ", mg.path)
	return nol.path, mapping, options, isMappingNil
}

func (nol *nonObliviousLocal) next(dest *graph.Node) (*graph.Node, []*graph.Node, int) {
	neighbors, neighIsNil := nol.network.GetNeighbors(nol.curr)
	//fmt.Println("Inside next - The neighbors are:", neighbors)
	//fmt.Println("Inside next")
	if neighIsNil {
		return nil, nil, -1
	}
	optimumNode := neighbors[0]
	tempNode := optimumNode
	choices := make([]*graph.Node, 1)
	choices[0] = optimumNode
	for _, node := range neighbors {
		if graph.IsEqual(node.ID, tempNode.ID) {
			continue
		}
		//////// This is the important part.
		if nol.network.GetLinkBetween(nol.curr, node).IsActive == false {
			//fmt.Println("Link is not active!")
			//fmt.Println("nol.curr is", nol.curr.ID, "node is", node.ID)
			continue
		}
		_, neighOfNeighIsNil := nol.network.GetNeighbors(node)
		if neighOfNeighIsNil {
			continue
		}
		if nol.network.Distance(node, dest, "hop") == nol.network.Distance(optimumNode, dest, "hop") {
			choices = append(choices, node)
		}
		if nol.network.Distance(node, dest, "hop") < nol.network.Distance(optimumNode, dest, "hop") {
			optimumNode = node
			choices = make([]*graph.Node, 1)
			choices[0] = optimumNode
		}
	}
	if len(choices) == 1 {
		return optimumNode, choices, 0
	} else {
		options := make([]*graph.Node, 1)
		//fmt.Println("Choices", len(choices))
		for _, node := range choices {
			link := nol.network.GetLinkBetween(nol.curr, node)
			if link.IsActivated() == true {
				options = append(options, node)
				//linkToPrune := make([]*graph.Link, 1)
				//linkToPrune[0] = link
				//graph.Prune(linkToPrune)
				//fmt.Println("This is pruned", link.ID)
				//return node, choices
			}
		}
		//if len(options) > 1 {
		//	r := 0
		//	for r == 0 {
		//fmt.Println("Gir Eladim!")
		//		r = rand.Intn(len(options))
		//	}
		//	return options[r], choices, r
		//}
		//r := 0
		//for r == 0 {
		//fmt.Println("Gir Eladim!")
		r := rand.Intn(len(choices))
		//}
		//link := mg.network.GetLinkBetween(mg.curr, choices[0])
		//linkToPrune := make([]*graph.Link, 1)
		//linkToPrune[0] = link
		//graph.Prune(linkToPrune)
		//fmt.Println("This is pruned", link.ID)
		return choices[r], choices, r
	}
}

func (nol *nonObliviousLocal) add(n *graph.Node) {
	//fmt.Println(nol.path)
	if nol.path[0].Memory == 0 {
		nol.path[0] = n
		//copy(nol.path[0], n)
		//fmt.Println("HERE!!!")
	} else {
		nol.path = append(nol.path, n)
	}
	//fmt.Println("Inside add. Input:", n)
	//fmt.Println("Inside add. Output:", nol.path)
}
