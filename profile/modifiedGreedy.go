package profile

import (
	"fmt"
	"math/rand"
	"time"

	"example.com/config"
	"example.com/graph"
	"example.com/path"

	"example.com/request"

	"example.com/quantum"
)

type ModifiedGreedyProfile struct {
	Network          *graph.Grid
	isFinished       bool
	hasRecovery      bool
	RunTime          int
	PriorityLen      int
	LinksWaitingTime []float64
	pathAlgorithm    string
}

func (mgp *ModifiedGreedyProfile) Build(topology string) {
	mgp.RunTime = 0
	mgp.PriorityLen = 3
	mgp.LinksWaitingTime = make([]float64, 0)
	mgp.hasRecovery = config.GetConfig().GetHasRecovery()
	mgp.pathAlgorithm = path.MODIFIED_GREEDY
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
		mgp.Network = grid
		//} else if topology == graph.RING {
		//ring := new(graph.Ring)
		//ring.Build()
		//mgp.Network = ring
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

func (mgp *ModifiedGreedyProfile) GenRequests(ignoreLeftOvers bool) []*request.Request {
	numRequests := config.GetConfig().GetNumRequests()
	var priority []int
	priority = make([]int, numRequests)
	rand.Seed(time.Now().UTC().UnixNano())
	// Priority for the requests
	for i := 0; i < numRequests; i++ {
		//priority[i] = 1
		priority[i] = rand.Intn(mgp.PriorityLen) + 1
	}
	//fmt.Println(priority)
	ids := mgp.Network.GetNodeIDs()
	reqs, err := request.RG(numRequests, ids, priority, mgp.Network.GetType(), mgp.RunTime)
	if err != nil {
		fmt.Println("Profile genRequests: Error in request generation!", err)
		return nil
	}
	//fmt.Println("Inside profile.GenRequests, behind path.PF")
	path.PF(mgp.Network, reqs, mgp.pathAlgorithm, ignoreLeftOvers)
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

////////////////// TODO: The requests should remain the same through iterations.
////////////////// Take them out.
func (mgp *ModifiedGreedyProfile) Run(reqs []*request.Request, maxItr int) {
	links := mgp.Network.GetLinks()
	numReached := 0
	isOpportunistic := config.GetConfig().GetIsOpportunistic()
	itrCntr := 0
	//var cntr int
	whichPath := make([]int, len(reqs))
	if !isOpportunistic {
		//fmt.Println("in NOPP!")
		mgp.isFinished = false
		for !mgp.isFinished {
			itrCntr++
			//numReached = 0
			if itrCntr == maxItr {
				break
			}
			linksWaiting := make([]float64, 0)
			///////////////////////////////// Check the following commented isReady.
			//isReady := true
			mgp.RunTime++
			// EG() also handles lifetimes.
			quantum.EG(links, mgp.RunTime)
			/*for rr, req := range reqs {
				fmt.Println("req number is", rr, "req source is", req.Src, "req dest is", req.Dest)
				fmt.Println("PositionID is", req.PositionID, "position is", req.Position)
				for mm := range req.Paths {
					fmt.Println("Path number", mm)
					for nn := range req.Paths[mm] {
						fmt.Println("The path for request:", req.Paths[mm][nn].ID)
					}
				}
			}*/
			if !mgp.hasRecovery {
				numReached, whichPath, linksWaiting = noRecoveryRun(mgp.Network, reqs, whichPath, numReached, mgp.RunTime, false)
			}
			mgp.LinksWaitingTime = append(mgp.LinksWaitingTime, linksWaiting...)
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				mgp.isFinished = true
			}
		}
	} else {
		//fmt.Println("in OPP!")
		mgp.isFinished = false
		numReached = 0
		for !mgp.isFinished {
			itrCntr++
			if itrCntr == maxItr {
				break
			}
			linksWaiting := make([]float64, 0)
			//numReached = 0
			//k := config.GetConfig().GetOpportunismDegree()
			//isReady := true
			mgp.RunTime++
			quantum.EG(links, mgp.RunTime)
			/*for rr, req := range reqs {
				fmt.Println("req number is", rr, "req source is", req.Src, "req dest is", req.Dest)
				fmt.Println("PositionID is", req.PositionID, "position is", req.Position)
				for mm := range req.Paths {
					fmt.Println("Path number", mm)
					for nn := range req.Paths[mm] {
						fmt.Println("The path for request:", req.Paths[mm][nn].ID)
					}
				}
			}*/
			if !mgp.hasRecovery {
				numReached, whichPath, linksWaiting = noRecoveryRunOPP(mgp.Network, reqs, whichPath, numReached, mgp.RunTime, false)
			}
			mgp.LinksWaitingTime = append(mgp.LinksWaitingTime, linksWaiting...)
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				mgp.isFinished = true
			}
		}
	}
}

func (mgp *ModifiedGreedyProfile) Stop() {
	mgp.isFinished = true
}

func (mgp *ModifiedGreedyProfile) Clear() {
	mgp.isFinished = false
	mgp.RunTime = 0
	mgp.LinksWaitingTime = make([]float64, 0)
	mgp.Network.Clear()
}

func (mgp *ModifiedGreedyProfile) GetNetwork() graph.Topology {
	return mgp.Network
}

func (mgp *ModifiedGreedyProfile) GetRunTime() int {
	return mgp.RunTime
}

func (mgp *ModifiedGreedyProfile) GetPriorityLen() int {
	return mgp.PriorityLen
}

func (mgp *ModifiedGreedyProfile) GetLinksWaitingTime() []float64 {
	return mgp.LinksWaitingTime
}

func (mgp *ModifiedGreedyProfile) GetHasRecovery() bool {
	return mgp.hasRecovery
}

func (mgp *ModifiedGreedyProfile) GetPathAlgorithm() string {
	return mgp.pathAlgorithm
}
