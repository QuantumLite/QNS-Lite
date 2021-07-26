package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/quantum"
	"example.com/request"
)

type NonObliviousLocal struct {
	Network          *graph.Grid
	isFinished       bool
	hasRecovery      bool
	RunTime          int
	PriorityLen      int
	LinksWaitingTime []float64
	pathAlgorithm    string
}

func (nol *NonObliviousLocal) Build(topology string) {
	nol.RunTime = 0
	nol.hasRecovery = config.GetConfig().GetHasRecovery()
	nol.pathAlgorithm = path.NONOBLIVIOUS_LOCAL
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
		nol.Network = grid
		//} else if topology == graph.RING {
		//ring := new(graph.Ring)
		//ring.Build()
		//nol.Network = ring
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

func (nol *NonObliviousLocal) GenRequests(ignoreLeftOvers bool) []*request.Request {
	numRequests := config.GetConfig().GetNumRequests()
	var priority []int
	priority = make([]int, numRequests)
	// Priority for the requests
	for i := 0; i < numRequests; i++ {
		priority[i] = 1
	}
	ids := nol.Network.GetNodeIDs()
	reqs, err := request.RG(numRequests, ids, priority, nol.Network.GetType(), nol.RunTime)
	if err != nil {
		fmt.Println("Profile genRequests: Error in request generation!", err)
		return nil
	}
	//fmt.Println("Inside profile.GenRequests, behind path.PF")
	path.PF(nol.Network, reqs, nol.pathAlgorithm, ignoreLeftOvers)
	//fmt.Println("Inside profile.GenRequests, after path.PF")

	/*for _, req := range reqs {
		n1 := req.Src
		n2 := req.Dest
		fmt.Println(*n1)
		fmt.Println(*n2)
		fmt.Println(len(req.Paths[0]))
		lenn := len(req.Paths)
		for i := 0; i <= lenn-1; i++ {
			for _, nodede := range req.Paths[i] {
				fmt.Println("The next node for path", i+1)
				fmt.Println(*nodede)
			}
		}
	}*/
	return reqs
}

func (nol *NonObliviousLocal) Run(reqs []*request.Request, maxItr int) {
	links := nol.Network.GetLinks()
	numReached := 0
	isOpportunistic := config.GetConfig().GetIsOpportunistic()
	itrCntr := 0
	//var cntr int
	whichPath := make([]int, len(reqs))

	if !isOpportunistic {
		nol.isFinished = false
		for !nol.isFinished {
			//fmt.Println("NOPP not finished yet.")
			itrCntr++
			//numReached = 0
			if itrCntr == maxItr {
				break
			}
			linksWaiting := make([]float64, 0)
			nol.RunTime++
			request.ClearReqPaths(reqs)
			// EG() also handles lifetimes.
			quantum.EG(links, nol.RunTime)
			//fmt.Println("Before path.PF in nonoblivious local.")
			pathlessReqs := make([]*request.Request, 0)
			for _, req := range reqs {
				if req.HasReached {
					continue
				}
				pathlessReqs = append(pathlessReqs, req)
			}
			//path.PF(nol.Network, reqs, "nonoblivious local", true)
			path.PF(nol.Network, pathlessReqs, "nonoblivious local", true)
			for _, req := range reqs {
				req.PositionID = req.Src.ID
				/*fmt.Println("req number is", rr, "req source is", req.Src, "req dest is", req.Dest)
				fmt.Println("PositionID is", req.PositionID, "position is", req.Position)
				for mm := range req.Paths {
					fmt.Println("Path number", mm)
					for nn := range req.Paths[mm] {
						fmt.Println("The path for request:", req.Paths[mm][nn].ID)
					}
				}*/
			}
			//fmt.Println("NOPP Found paths")
			whichPath = make([]int, len(reqs))
			if !nol.hasRecovery {
				numReached, _, linksWaiting = noRecoveryRun(nol.Network, reqs, whichPath, numReached, nol.RunTime, true)
			}
			nol.LinksWaitingTime = append(nol.LinksWaitingTime, linksWaiting...)
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				nol.isFinished = true
			}
		}
	} else {
		nol.isFinished = false
		numReached = 0
		//fmt.Println("Inside OPP. isFinished is:", nol.isFinished)
		for !nol.isFinished {
			//fmt.Println("OPP not finished yet.")
			itrCntr++
			if itrCntr == maxItr {
				break
			}
			linksWaiting := make([]float64, 0)
			//numReached = 0
			//k := config.GetConfig().GetOpportunismDegree()
			//isReady := true
			nol.RunTime++
			request.ClearReqPaths(reqs)
			quantum.EG(links, nol.RunTime)
			pathlessReqs := make([]*request.Request, 0)
			for _, req := range reqs {
				if req.HasReached {
					continue
				}
				pathlessReqs = append(pathlessReqs, req)
			}
			//path.PF(nol.Network, reqs, "nonoblivious local", true)
			path.PF(nol.Network, pathlessReqs, "nonoblivious local", true)
			for _, req := range reqs {
				req.PositionID = req.Src.ID
				/*fmt.Println("req number is", rr, "req source is", req.Src, "req dest is", req.Dest)
				fmt.Println("PositionID is", req.PositionID, "position is", req.Position)
				for mm := range req.Paths {
					fmt.Println("Path number", mm)
					for nn := range req.Paths[mm] {
						fmt.Println("The path for request:", req.Paths[mm][nn].ID)
					}
				}*/
			}
			//for n, req := range reqs {
			//	for m := 0; m < len(req.Paths); m++ {
			//		for _, node := range req.Paths[m] {
			//			fmt.Println("req is", n, "path is", m, "Node is", node.ID)
			//		}
			//	}
			//}
			whichPath = make([]int, len(reqs))
			if !nol.hasRecovery {
				numReached, _, linksWaiting = noRecoveryRunOPP(nol.Network, reqs, whichPath, numReached, nol.RunTime, true)
			}
			nol.LinksWaitingTime = append(nol.LinksWaitingTime, linksWaiting...)
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				nol.isFinished = true
			}
		}
	}
}

func (nol *NonObliviousLocal) Stop() {
	nol.isFinished = true
}

func (nol *NonObliviousLocal) Clear() {
	nol.isFinished = false
	//fmt.Println("Cleared! isFinished is", nol.isFinished)
	nol.RunTime = 0
	nol.LinksWaitingTime = make([]float64, 0)
	nol.Network.Clear()
}

func (nol *NonObliviousLocal) GetNetwork() graph.Topology {
	return nol.Network
}

func (nol *NonObliviousLocal) GetRunTime() int {
	return nol.RunTime
}

func (nol *NonObliviousLocal) GetLinksWaitingTime() []float64 {
	return nol.LinksWaitingTime
}

func (nol *NonObliviousLocal) GetPriorityLen() int {
	return nol.PriorityLen
}

func (nol *NonObliviousLocal) GetHasRecovery() bool {
	return nol.hasRecovery
}

func (nol *NonObliviousLocal) GetPathAlgorithm() string {
	return nol.pathAlgorithm
}
