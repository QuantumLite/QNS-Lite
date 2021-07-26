package log

import (
	"fmt"

	"example.com/graph"
	"example.com/request"
)

func PrintPaths(reqs []*request.Request) {
	requests := make([]*request.Request, 1)
	for reqIndex, req := range reqs {
		requests[0] = req
		PrintReqs(requests)
		for i := 0; i <= len(req.Paths)-1; i++ {
			for _, node := range req.Paths[i] {
				fmt.Println("The next node for request", reqIndex, ", path", i)
				fmt.Println(*node)
			}
		}
	}
}

func PrintRecoveryPaths(reqs []*request.Request) {
	requests := make([]*request.Request, 1)
	for reqIndex, req := range reqs {
		requests[0] = req
		PrintReqs(requests)
		for pathIndex, _ := range req.Paths {
			for cursor, _ := range req.RecoveryPaths[pathIndex] {
				for recIndex, recoveryPath := range req.RecoveryPaths[pathIndex][cursor] {
					for _, node := range recoveryPath {
						fmt.Println("The next recovery node for request", reqIndex, ", path", pathIndex, ", cursor", cursor, "recovery Index", recIndex)
						fmt.Println(*node)
					}
				}
			}
		}
	}
}

func PrintReqs(reqs []*request.Request) {
	for index, req := range reqs {
		n1 := req.Src
		n2 := req.Dest
		fmt.Println("Request number", index)
		fmt.Println("Source is:")
		fmt.Println(*n1)
		fmt.Println("Destination is:")
		fmt.Println(*n2)
	}
}

func PrintLinks(links []*graph.Link) {
	fmt.Println("Log - Printing the link IDs.")
	for i, link := range links {
		fmt.Println("Link number:", i, "The link ID is:", link.ID)
	}
}
