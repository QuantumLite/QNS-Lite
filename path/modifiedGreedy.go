package path

import (
	"math/rand"
	"time"

	"example.com/graph"
	"example.com/utils"
)

type modifiedGreedy struct {
	network    graph.Topology
	path       Path
	isFinished bool
	curr       *graph.Node
	//src        *graph.Node
	//dest       *graph.Node
}

func (mg *modifiedGreedy) Build(network graph.Topology) {
	mg.path = make([]*graph.Node, 1)
	mg.path[0] = new(graph.Node)
	mg.curr = new(graph.Node)
	mg.isFinished = false
	mg.network = network
	rand.Seed(time.Now().UTC().UnixNano())
	//mg.path[0] = mg.src
	//mg.curr = mg.src
}

// Clear flushes the path finder after it has found a path for a request.
func (mg *modifiedGreedy) Clear() {
	mg.path = nil
	mg.isFinished = false
	mg.path = make([]*graph.Node, 1)
	mg.path[0] = new(graph.Node)
	//mg.curr = nil
}

func (mg *modifiedGreedy) Find(src, dest *graph.Node) (Path, []int, []int, bool) {
	mg.curr = src
	mg.add(src)
	cntr := 0
	//mapping := make(map[*graph.Node]int)
	mapping := make([]int, 1)
	isMappingNil := true
	mappingCursor := 0
	options := make([]int, 1)
	counter := 0
	anotherCounter := 0
	for !mg.curr.IsEqual(dest) {
		if cntr >= mg.network.GetSize()*mg.network.GetSize() {
			//fmt.Println("Inside Find. Counter overflow.")
			return nil, nil, nil, true
		}
		cntr = cntr + 1
		next := mg.path[len(mg.path)-1]
		var choices []*graph.Node
		check := true
		for check {
			anotherCounter++
			if anotherCounter >= mg.network.GetSize()*mg.network.GetSize() {
				//fmt.Println("Inside Find. anotherCounter overflow.")
				return nil, nil, nil, true
			}
			if len(mg.path) < 2 {
				temp, tempChoices, _ := mg.next(dest)
				next = temp
				choices = tempChoices
				break
			}
			temp, tempChoices, _ := mg.next(dest)
			next = temp
			choices = tempChoices
			if next == nil {
				break
			}
			if graph.IsEqual(next.ID, mg.path[len(mg.path)-2].ID) {
				//fmt.Println("Find modifiedGreedy - AHAAAAAAAAAAAAAAAAAAAAAAAAA!")
				//fmt.Println("next.ID is", next.ID)
				//fmt.Println("The node before curr is", mg.path[len(mg.path)-2].ID)
				//fmt.Println("curr is", mg.curr)
			}
			if !graph.IsEqual(next.ID, mg.path[len(mg.path)-2].ID) {
				////////////////////////////////////// Maybe I can do some pruning here.
				check = false
			}
		}

		if len(choices) > 1 {
			if len(mg.path) >= 2 {
				//fmt.Println("Check here - 55555555555555555555555555555555555555555555555555555555")
				//for _, choiceNode := range choices {
				//	fmt.Println("node is", choiceNode.ID)
				//}
				for choiceIndex, choiceNode := range choices {
					if graph.IsEqual(choiceNode.ID, mg.path[len(mg.path)-2].ID) {
						//fmt.Println("poxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", choiceNode.ID)
						choices = utils.RemoveNode(choices, choiceIndex)
					}
				}
				//for _, choiceNode := range choices {
				//	fmt.Println("after - node is", choiceNode.ID)
				//}
			}
		}
		//fmt.Println("next is", next.ID, "last is", mg.path[len(mg.path)-1])
		if len(choices) > 1 {
			isMappingNil = false
			if mappingCursor == 0 {
				mapping[mappingCursor] = len(mg.path) - 1
			} else {
				mapping = append(mapping, len(mg.path)-1)
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
		mg.add(next)
		mg.curr = mg.path[len(mg.path)-1]
	}
	//fmt.Println("Found Path - inside find: ", mg.path)
	//for _, nodede := range mg.path {
	//	fmt.Println(nodede.ID)
	//}
	return mg.path, mapping, options, isMappingNil
}

func (mg *modifiedGreedy) next(dest *graph.Node) (*graph.Node, []*graph.Node, int) {
	neighbors, neighIsNil := mg.network.GetNeighbors(mg.curr)
	//fmt.Println("Inside next - The neighbors are:", neighbors)
	if neighIsNil {
		//fmt.Println("Nil neighbors!!!!!!!!!!!!!!!!")
		return nil, nil, -1
	}
	optimumNode := neighbors[0]
	tempNode := optimumNode
	optimumNeighOfNeigh, optimumNeighOfNeighIsNil := mg.network.GetNeighbors(optimumNode)
	choices := make([]*graph.Node, 0)
	if !optimumNeighOfNeighIsNil && len(optimumNeighOfNeigh) > 1 {
		choices = append(choices, optimumNode)
	}
	//choices[0] = optimumNode
	for _, node := range neighbors {
		if graph.IsEqual(node.ID, tempNode.ID) {
			//optimumNeighOfNeigh, optimumNeighOfNeighIsNil := mg.network.GetNeighbors(node)
			//if optimumNeighOfNeighIsNil || len(optimumNeighOfNeigh) == 1 {
			//	choices = make([]*graph.Node, 0)
			//}
			continue
		}
		neighOfNeigh, neighOfNeighIsNil := mg.network.GetNeighbors(node)
		if neighOfNeighIsNil || len(neighOfNeigh) == 1 {
			//fmt.Println("modifiedGreedy.go path neighOfNeigh: OOOOOOOOOOOOOOOOOOOOOOOPS!!!!!!!!")
			continue
		}
		/////////////////// Added this to solve one strange behavior of findRecoveryPaths.
		/////////////////// IMPORTANT!!!!!!!!!!!!!
		if mg.network.Distance(node, dest, "hop") > mg.network.Distance(mg.curr, dest, "hop") {
			continue
		}
		if mg.network.Distance(node, dest, "hop") == mg.network.Distance(optimumNode, dest, "hop") {
			choices = append(choices, node)
		}
		if mg.network.Distance(node, dest, "hop") < mg.network.Distance(optimumNode, dest, "hop") {
			optimumNode = node
			choices = make([]*graph.Node, 1)
			choices[0] = optimumNode
		}
	}
	if len(choices) == 0 {
		return nil, nil, -1
	}
	if len(choices) == 1 {
		return optimumNode, choices, 0
	} else {
		options := make([]*graph.Node, 1)
		//fmt.Println("Choices", len(choices))
		for _, node := range choices {
			link := mg.network.GetLinkBetween(mg.curr, node)
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

func (mg *modifiedGreedy) add(n *graph.Node) {
	//fmt.Println(mg.path)
	if mg.path[0].Memory == 0 {
		mg.path[0] = n
		//copy(mg.path[0], n)
		//fmt.Println("HERE!!!")
	} else {
		mg.path = append(mg.path, n)
	}
	//fmt.Println("Inside add. Input:", n)
	//fmt.Println("Inside add. Output:", mg.path)
}
