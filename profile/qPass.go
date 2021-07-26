package profile

import (
	"fmt"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/quantum"
	"example.com/request"
)

type QPass struct {
	Network          *graph.Grid
	isFinished       bool
	hasRecovery      bool
	RunTime          int
	PriorityLen      int
	LinksWaitingTime []float64
	pathAlgorithm    string
}

func (qpass *QPass) Build(topology string) {
	qpass.RunTime = 0
	qpass.hasRecovery = config.GetConfig().GetHasRecovery()
	qpass.pathAlgorithm = path.MODIFIED_GREEDY
	if topology == graph.GRID {
		grid := new(graph.Grid)
		grid.Build()
		qpass.Network = grid
		//} else if topology == graph.RING {
		//ring := new(graph.Ring)
		//ring.Build()
		//qpass.Network = ring
	} else {
		fmt.Println("Profile: Caution! The topology is not implemented.")
	}
}

func (qpass *QPass) GenRequests(ignoreLeftOvers bool) []*request.Request {
	numRequests := config.GetConfig().GetNumRequests()
	var priority []int
	priority = make([]int, numRequests)
	// Priority for the requests
	for i := 0; i < numRequests; i++ {
		priority[i] = 1
	}
	ids := qpass.Network.GetNodeIDs()
	reqs, err := request.RG(numRequests, ids, priority, qpass.Network.GetType(), qpass.RunTime)
	if err != nil {
		fmt.Println("Profile genRequests: Error in request generation!", err)
		return nil
	}
	//fmt.Println("Inside profile.GenRequests, behind path.PF")
	path.PF(qpass.Network, reqs, qpass.pathAlgorithm, ignoreLeftOvers)
	findRecoveryPaths(reqs, qpass.Network)
	//log.PrintPaths(reqs)
	//log.PrintRecoveryPaths(reqs)

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

func (qpass *QPass) Run(reqs []*request.Request, maxItr int) {
	links := qpass.Network.GetLinks()
	////// Uncomment!!!
	numReached := 0
	isOpportunistic := config.GetConfig().GetIsOpportunistic()
	itrCntr := 0
	//var cntr int

	////// Uncomment!!!
	whichPath := make([]int, len(reqs))
	if !isOpportunistic {
		qpass.isFinished = false
		for !qpass.isFinished {
			itrCntr++
			//numReached = 0
			if itrCntr == maxItr {
				break
			}
			linksWaiting := make([]float64, 0)
			///////////////////////////////// Check the following commented isReady.
			//isReady := true
			qpass.RunTime++
			// EG() also handles lifetimes.
			quantum.EG(links, qpass.RunTime)
			if !qpass.hasRecovery {
				numReached, whichPath, linksWaiting = RecoveryRun(qpass.Network, reqs, whichPath, numReached, qpass.RunTime, false)
			}
			qpass.LinksWaitingTime = append(qpass.LinksWaitingTime, linksWaiting...)
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)

			//////// Uncomment!!!
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				qpass.isFinished = true
			}
		}
	} else {
		qpass.isFinished = false
		//// Uncomment!!!!
		numReached = 0
		for !qpass.isFinished {
			itrCntr++
			if itrCntr == maxItr {
				break
			}
			linksWaiting := make([]float64, 0)
			//numReached = 0
			//k := config.GetConfig().GetOpportunismDegree()
			//isReady := true
			qpass.RunTime++
			quantum.EG(links, qpass.RunTime)
			if !qpass.hasRecovery {
				numReached, whichPath, linksWaiting = recoveryRunOPP(qpass.Network, reqs, whichPath, numReached, qpass.RunTime, false)
			}
			qpass.LinksWaitingTime = append(qpass.LinksWaitingTime, linksWaiting...)
			//fmt.Println("Number of reached::::::::::::::::::::::", numReached)

			/////// Uncomment!!!
			if numReached == len(reqs) {
				//fmt.Println("REACHED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				qpass.isFinished = true
			}
		}
	}
}

func (qpass *QPass) Stop() {
	qpass.isFinished = true
}

func (qpass *QPass) Clear() {
	qpass.isFinished = false
	qpass.RunTime = 0
	qpass.LinksWaitingTime = make([]float64, 0)
	qpass.Network.Clear()
}

func (qpass *QPass) GetNetwork() graph.Topology {
	return qpass.Network
}

func (qpass *QPass) GetRunTime() int {
	return qpass.RunTime
}

func (qpass *QPass) GetLinksWaitingTime() []float64 {
	return qpass.LinksWaitingTime
}

func (qpass *QPass) GetPriorityLen() int {
	return qpass.PriorityLen
}

func (qpass *QPass) GetHasRecovery() bool {
	return qpass.hasRecovery
}

func (qpass *QPass) GetPathAlgorithm() string {
	return qpass.pathAlgorithm
}
