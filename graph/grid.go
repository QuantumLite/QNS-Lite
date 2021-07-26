package graph

import (
	"fmt"
	"math"

	"example.com/config"
)

// The constant variables
const (
	GRID_NODE_ID_DIM = 2
	GRID_LINK_ID_DIM = 3
)

type Grid struct {
	gridSize int
	Nodes    [][]*Node
	Links    [][][]*Link
	Type     string
}

func (grid *Grid) Build() {
	grid.gridSize = config.GetConfig().GetSize()
	grid.Type = "grid"
	grid.generateNodes()
	//fmt.Println("Node IDs", grid.nodes)
	grid.generateLinks()
}

func (grid *Grid) Clear() {
	links := grid.GetLinks()
	for _, link := range links {
		link.GenerationTime = 0
		link.ConsumptionTime = 0
		link.IsActive = false
		link.IsReserved = false
		link.Reservation = -1
	}
}

func (grid *Grid) generateNodes() {
	id := make([]int, 2)
	grid.Nodes = make([][]*Node, grid.gridSize)
	for i := 0; i < grid.gridSize; i++ {
		grid.Nodes[i] = make([]*Node, grid.gridSize)
	}
	var node *Node
	for i := 0; i < grid.gridSize; i++ {
		for j := 0; j < grid.gridSize; j++ {
			id[0] = i
			id[1] = j
			//fmt.Println(id)
			grid.Nodes[i][j] = new(Node)
			//////// IMPORTANT: CHECK THIS!
			grid.Nodes[i][j].ID = make([]int, 2)
			node = MakeNode(id, config.GetConfig().GetMemory())
			copy(grid.Nodes[i][j].ID, node.ID)
			grid.Nodes[i][j].Memory = node.Memory
			//fmt.Println(MakeNode(id, config.GetConfig().GetMemory()))
			//fmt.Println("Node ID", grid.nodes[i][2].ID)
		}
	}
}

func (grid *Grid) generateLinks() {
	id := make([]int, 3)
	//lifetime := config.GetConfig().GetLifetime()
	grid.Links = make([][][]*Link, grid.gridSize)
	var link *Link
	for i := 0; i < grid.gridSize; i++ {
		grid.Links[i] = make([][]*Link, 2)
		for j := 0; j < 2; j++ {
			grid.Links[i][j] = make([]*Link, grid.gridSize-1)
			for k := 0; k < grid.gridSize-1; k++ {
				grid.Links[i][j][k] = new(Link)
				grid.Links[i][j][k].ID = make([]int, 3)
				id[0] = i
				id[1] = j
				id[2] = k
				link = MakeLink(id, 0, false)
				copy(grid.Links[i][j][k].ID, link.ID)
				grid.Links[i][j][k].Age = link.Age
				grid.Links[i][j][k].GenerationTime = 0
				grid.Links[i][j][k].ConsumptionTime = 0
				grid.Links[i][j][k].IsActive = link.IsActive
				grid.Links[i][j][k].IsPruned = link.IsPruned
				grid.Links[i][j][k].Reservation = link.Reservation
				grid.Links[i][j][k].IsReserved = link.IsReserved
				//grid.links[i][j][k] = MakeLink(id, lifetime, false)
			}
		}
	}
}

// TODO: Complete this!!!!
func (grid *Grid) GetLinks() []*Link {
	links := make([]*Link, grid.gridSize*2*(grid.gridSize-1))
	for i := 0; i < grid.gridSize; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < grid.gridSize-1; k++ {
				links[i*2*(grid.gridSize-1)+j*(grid.gridSize-1)+k] = grid.Links[i][j][k]
			}
		}
	}
	return links
}

//func GetNodes(grid *Grid) [][]*Node {
//	return grid.nodes
//}

func (grid *Grid) GetSize() int {
	return grid.gridSize
}

func (grid *Grid) GetNodeIDs() [][]int {
	ids := make([][]int, grid.gridSize*grid.gridSize)
	for i := 0; i < grid.gridSize*grid.gridSize; i++ {
		ids[i] = make([]int, 2)
	}
	//fmt.Println(ids[0])
	for i := 0; i < grid.gridSize; i++ {
		for j := 0; j < grid.gridSize; j++ {
			ids[i*grid.gridSize+j] = grid.Nodes[i][j].ID
		}
	}
	return ids
}

func (grid *Grid) GetNeighbors(node *Node) ([]*Node, bool) {
	neighbors := make([]*Node, 1)
	neighbors[0] = new(Node)
	//neighbors[1] = new(Node)
	//fmt.Println(neighbors)
	for _, n := range neighbors {
		n.ID = make([]int, 2)
	}
	bottomLeft := []int{0, 0}
	bottomRight := []int{0, grid.gridSize - 1}
	topLeft := []int{grid.gridSize - 1, 0}
	topRight := []int{grid.gridSize - 1, grid.gridSize - 1}
	isNil := true
	if IsEqual(node.ID, bottomLeft) {
		if grid.GetLinkBetween(node, grid.Nodes[1][0]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[1][0]
			} else {
				neighbors = append(neighbors, grid.Nodes[1][0])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[0][1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[0][1]
			} else {
				neighbors = append(neighbors, grid.Nodes[0][1])
			}
		}
		//neighbors[0] = grid.Nodes[1][0]
		//neighbors[1] = grid.Nodes[0][1]
		return neighbors, isNil
	}
	if IsEqual(node.ID, bottomRight) {
		if grid.GetLinkBetween(node, grid.Nodes[1][grid.gridSize-1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[1][grid.gridSize-1]
			} else {
				neighbors = append(neighbors, grid.Nodes[1][grid.gridSize-1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[0][grid.gridSize-2]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[0][grid.gridSize-2]
			} else {
				neighbors = append(neighbors, grid.Nodes[0][grid.gridSize-2])
			}
		}
		//neighbors[0] = grid.Nodes[1][grid.gridSize-1]
		//neighbors[1] = grid.Nodes[0][grid.gridSize-2]
		return neighbors, isNil
	}
	if IsEqual(node.ID, topLeft) {
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-1][1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-1][1]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-1][1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-2][0]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-2][0]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-2][0])
			}
		}
		//neighbors[0] = grid.Nodes[grid.gridSize-1][1]
		//neighbors[1] = grid.Nodes[grid.gridSize-2][0]
		return neighbors, isNil
	}
	if IsEqual(node.ID, topRight) {
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-2][grid.gridSize-1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-2][grid.gridSize-1]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-2][grid.gridSize-1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-1][grid.gridSize-2]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-1][grid.gridSize-2]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-1][grid.gridSize-2])
			}
		}
		//neighbors[0] = grid.Nodes[grid.gridSize-2][grid.gridSize-1]
		//neighbors[1] = grid.Nodes[grid.gridSize-1][grid.gridSize-2]
		return neighbors, isNil
	}
	if node.ID[0] == 0 {
		if grid.GetLinkBetween(node, grid.Nodes[0][node.ID[1]-1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[0][node.ID[1]-1]
			} else {
				neighbors = append(neighbors, grid.Nodes[0][node.ID[1]-1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[0][node.ID[1]+1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[0][node.ID[1]+1]
			} else {
				neighbors = append(neighbors, grid.Nodes[0][node.ID[1]+1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[1][node.ID[1]]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[1][node.ID[1]]
			} else {
				neighbors = append(neighbors, grid.Nodes[1][node.ID[1]])
			}
		}
		//neighbors[0] = grid.Nodes[0][node.ID[1]-1]
		//neighbors[1] = grid.Nodes[0][node.ID[1]+1]
		//neighbors = append(neighbors, grid.Nodes[1][node.ID[1]])
		return neighbors, isNil
	}
	if node.ID[0] == grid.gridSize-1 {
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-1][node.ID[1]-1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-1][node.ID[1]-1]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-1][node.ID[1]-1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-1][node.ID[1]+1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-1][node.ID[1]+1]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-1][node.ID[1]+1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[grid.gridSize-2][node.ID[1]]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[grid.gridSize-2][node.ID[1]]
			} else {
				neighbors = append(neighbors, grid.Nodes[grid.gridSize-2][node.ID[1]])
			}
		}
		//neighbors[0] = grid.Nodes[grid.gridSize-1][node.ID[1]-1]
		//neighbors[1] = grid.Nodes[grid.gridSize-1][node.ID[1]+1]
		//neighbors = append(neighbors, grid.Nodes[grid.gridSize-2][node.ID[1]])
		return neighbors, isNil
	}
	if node.ID[1] == 0 {
		if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]-1][0]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[node.ID[0]-1][0]
			} else {
				neighbors = append(neighbors, grid.Nodes[node.ID[0]-1][0])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]+1][0]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[node.ID[0]+1][0]
			} else {
				neighbors = append(neighbors, grid.Nodes[node.ID[0]+1][0])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]][1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[node.ID[0]][1]
			} else {
				neighbors = append(neighbors, grid.Nodes[node.ID[0]][1])
			}
		}
		//neighbors[0] = grid.Nodes[node.ID[0]-1][0]
		//neighbors[1] = grid.Nodes[node.ID[0]+1][0]
		//neighbors = append(neighbors, grid.Nodes[node.ID[0]][1])
		return neighbors, isNil
	}
	if node.ID[1] == grid.gridSize-1 {
		if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]-1][grid.gridSize-1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[node.ID[0]-1][grid.gridSize-1]
			} else {
				neighbors = append(neighbors, grid.Nodes[node.ID[0]-1][grid.gridSize-1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]+1][grid.gridSize-1]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[node.ID[0]+1][grid.gridSize-1]
			} else {
				neighbors = append(neighbors, grid.Nodes[node.ID[0]+1][grid.gridSize-1])
			}
		}
		if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]][grid.gridSize-2]).IsPruned == false {
			isNil = false
			if neighbors[0].Memory == 0 {
				neighbors[0] = grid.Nodes[node.ID[0]][grid.gridSize-2]
			} else {
				neighbors = append(neighbors, grid.Nodes[node.ID[0]][grid.gridSize-2])
			}
		}
		//neighbors[0] = grid.Nodes[node.ID[0]-1][grid.gridSize-1]
		//neighbors[1] = grid.Nodes[node.ID[0]+1][grid.gridSize-1]
		//neighbors = append(neighbors, grid.Nodes[node.ID[0]][grid.gridSize-2])
		return neighbors, isNil
	}
	if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]-1][node.ID[1]]).IsPruned == false {
		isNil = false
		if neighbors[0].Memory == 0 {
			neighbors[0] = grid.Nodes[node.ID[0]-1][node.ID[1]]
		} else {
			neighbors = append(neighbors, grid.Nodes[node.ID[0]-1][node.ID[1]])
		}
	}
	if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]][node.ID[1]-1]).IsPruned == false {
		isNil = false
		if neighbors[0].Memory == 0 {
			neighbors[0] = grid.Nodes[node.ID[0]][node.ID[1]-1]
		} else {
			neighbors = append(neighbors, grid.Nodes[node.ID[0]][node.ID[1]-1])
		}
	}
	if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]][node.ID[1]+1]).IsPruned == false {
		isNil = false
		if neighbors[0].Memory == 0 {
			neighbors[0] = grid.Nodes[node.ID[0]][node.ID[1]+1]
		} else {
			neighbors = append(neighbors, grid.Nodes[node.ID[0]][node.ID[1]+1])
		}
	}
	if grid.GetLinkBetween(node, grid.Nodes[node.ID[0]+1][node.ID[1]]).IsPruned == false {
		isNil = false
		if neighbors[0].Memory == 0 {
			neighbors[0] = grid.Nodes[node.ID[0]+1][node.ID[1]]
		} else {
			neighbors = append(neighbors, grid.Nodes[node.ID[0]+1][node.ID[1]])
		}
	}
	//neighbors[0] = grid.Nodes[node.ID[0]-1][node.ID[1]]
	//neighbors[1] = grid.Nodes[node.ID[0]][node.ID[1]-1]
	//neighbors = append(neighbors, grid.Nodes[node.ID[0]][node.ID[1]+1])
	//neighbors = append(neighbors, grid.Nodes[node.ID[0]+1][node.ID[1]])
	return neighbors, isNil
}

//////////////////////////////////////////// TODO: Check this function!
func (grid *Grid) GetLinkBetween(n1, n2 *Node) *Link {
	id1 := n1.ID
	id2 := n2.ID
	var x, y, z int
	if id1[0] == id2[0] && id1[1] == id2[1] {
		fmt.Println("Inside GetLinkBetween: identical nodes.")
		return nil
	}
	if id1[0] == id2[0] {
		x = id1[0]
		y = 0
		if id1[1] <= id2[1] {
			z = id1[1]
		} else {
			z = id2[1]
		}
	} else {
		x = id1[1]
		y = 1
		if id1[0] <= id2[0] {
			z = id1[0]
		} else {
			z = id2[0]
		}
	}
	//fmt.Println("Inside get links", x, y, z)
	return grid.Links[x][y][z]
}

func (grid *Grid) Distance(n1, n2 *Node, measure string) int {
	if measure == HOP {
		return int(math.Abs(float64(n1.ID[0]-n2.ID[0])) + math.Abs(float64(n1.ID[1]-n2.ID[1])))
	}
	return -1
}

func (grid *Grid) GetType() string {
	return grid.Type
}
