package profile

import (
	"fmt"
	"math"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/quantum"
	"example.com/request"
)

const (
	MODIFIED_GREEDY    = "modified greedy"
	NONOBLIVIOUS_LOCAL = "nonoblivious local"
	QPASS              = "qpass"
)

type Profile interface {
	Build(topology string)
	Run(reqs []*request.Request, maxItr int)
	Stop()
	Clear()
	GenRequests(ignoreLeftOvers bool) []*request.Request
	GetRunTime() int
	GetLinksWaitingTime() []float64
	GetPriorityLen() int
	GetHasRecovery() bool
	GetNetwork() graph.Topology
	GetPathAlgorithm() string
}

func GenRequests(numRequests int, network graph.Topology, isMultiPath bool, algorithm string, ignoreLeftOvers bool) []*request.Request {
	var priority []int
	priority = make([]int, numRequests)
	// Priority for the requests
	for i := 0; i < numRequests; i++ {
		priority[i] = 1
	}
	ids := network.GetNodeIDs()
	reqs, err := request.RG(numRequests, ids, priority, network.GetType(), 1)
	if err != nil {
		fmt.Println("Profile genRequests: Error in request generation!", err)
		return nil
	}
	//fmt.Println("Inside profile.GenRequests, behind path.PF")
	path.PF(network, reqs, algorithm, ignoreLeftOvers)
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

func noRecoveryRun(network graph.Topology, reqs []*request.Request, whichPath []int, numReached int, runTime int, changeSrc bool) (int, []int, []float64) {
	//numReached := 0
	//fmt.Println("Hiaaa!")
	isReady := true
	//whichPath := make([]int, len(reqs))
	linksWaiting := make([]float64, 0)
	var reachedNum int
	var w float64
	var cntr int
	for reqNum, req := range reqs {
		//fmt.Println("Run - The req is: ", reqNum, " The path is: ", whichPath[reqNum])
		if req.HasReached {
			// Release the reserved links
			// Here, req.CanMove is used to release the links, and is set to false
			// to prevent extra work every time the request enters this if statement.
			if req.CanMove {
				//fmt.Println("Req ", reqNum, "Freed the resources.")
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
					//fmt.Println("Freeing link", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
				}
			}
			req.CanMove = false
			continue
		}
		if len(req.Paths[0]) == 0 {
			continue
		}
		cntr = 0
		if !req.CanMove {
			for which, _ := range req.Paths {
				//fmt.Println("Which is ", which, "req is ", reqNum)
				cntr = 0
				isReady = true
				for i := 1; i <= len(req.Paths[which])-1; i++ {
					link := network.GetLinkBetween(req.Paths[which][i], req.Paths[which][i-1])
					if link.IsReserved == false {
						//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
							isReady = isReady && link.IsActive
							cntr++
						} else {
							//fmt.Println("1--- The req is ", reqNum, " It is unfortunately reserved by: ", link.Reservation)
							//fmt.Println("link is", link.ID)
							isReady = false
						}
					}
					if cntr == 0 {
						isReady = false
					}
					if !isReady {
						///////////////// IMPORTANT!!! CHECK THIS!!!!!
						break
					}
					// Solve the isReady issue.
				}
				if isReady {
					//fmt.Println("Reservation success!")
					//fmt.Println("Run - The req is ", reqNum, "which is ", which)
					whichPath[reqNum] = which
					if changeSrc {
						req.Position = graph.FindPosition(req.PositionID, req.Paths[which])
					}
					break
				}
			}
		} else {
			//fmt.Println("Have already reserved!")
			////////////////////////////////////////// VERY IMPORTANT!!!!!!!!!!!!!!
			//for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
			for i := req.Position + 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				link := network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1])
				if link.IsReserved == false {
					//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
					isReady = isReady && link.IsActive
					cntr++
				} else {
					if link.Reservation == reqNum || link.Reservation == -1 {
						//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					} else {
						//fmt.Println("2--- The req is ", reqNum, " It is unfortunately reserved by: ", link.Reservation)
						//fmt.Println("link is", link.ID)
						isReady = false
					}
				}
				if cntr == 0 {
					//fmt.Println("zero counter!")
					isReady = false
				}
				if !isReady {
					//fmt.Println("Profile: Damn!")
					break
				}
				// Solve the isReady issue.
			}
		}
		//fmt.Println("Request", reqNum, isReady)
		if isReady {
			//fmt.Println("profile: It is ready!")
			// req.CanMove shows the fact that the request has previously reserved the
			// path, and is only trying to swap its way to the end.
			if !req.CanMove {
				//fmt.Println("Req ", reqNum, "is reserving for path ", whichPath[reqNum])
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
					//fmt.Println("Reserving link: ", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
				}
			}
			req.CanMove = true
			//fmt.Println("------------------------LENGTH IS: ", len(req.Paths[whichPath[reqNum]]), "WHICH IS: ", whichPath[reqNum], "LENGTH OF PATHS IS: ", len(req.Paths))
			reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, false)
			numReached += reachedNum
			if w > 0 {
				linksWaiting = append(linksWaiting, w)
			}
			//if req.HasReached {
			//fmt.Println("Req ", reqNum, " Has reached!")
			//}
		}
		isReady = true
		if changeSrc {
			for m := 1; m <= len(req.Paths[whichPath[reqNum]])-1; m++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).Reservation = -1
			}
			req.CanMove = false
		}
	}
	return numReached, whichPath, linksWaiting
}

func noRecoveryRunOPP(network graph.Topology, reqs []*request.Request, whichPath []int, numReached int, runTime int, changeSrc bool) (int, []int, []float64) {
	//fmt.Println("Hiaaa OPP!")
	isReady := true
	oppCntr := 0
	k := config.GetConfig().GetOpportunismDegree()
	//whichPath := make([]int, len(reqs))
	linksWaiting := make([]float64, 0)
	var reachedNum int
	var w float64
	var cntr int
	for reqNum, req := range reqs {
		//fmt.Println("Beginning - reqNum is", reqNum, "req.Position is", req.Position, "req.PositionID is", req.PositionID, "req.Src is", req.Src.ID, "req.Dest is", req.Dest.ID)
		//config.SetOpportunismDegree(req.Priority)
		//k = config.GetConfig().GetOpportunismDegree()
		oppCntr = 0
		if req.CanMove {
			for i := 1; i <= req.Position-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
			}
		}
		if req.HasReached {
			// Release the reserved links
			// Here, req.CanMove is used to release the links, and is set to false
			// to prevent extra work every time the request enters this if statement.
			if req.CanMove {
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
				}
			}
			req.CanMove = false
			continue
		}
		cntr = 0
		if len(req.Paths[0]) == 0 {
			continue
		}
		// req.Position starts from 1. Check this!!!!!!!!!!!!!!!!!!!!!!!!!
		if !req.CanMove {
			//fmt.Println("Can't move OPP! reqNum is", reqNum)
			for which, _ := range req.Paths {
				/////////////////// Fill in here!
				//fmt.Println("Which is ", which, "req is ", reqNum)
				oppCntr = 0
				pos := graph.FindPosition(req.PositionID, req.Paths[which])
				//fmt.Println("Which is ", which, "pos is ", pos, "len(req.Paths[which])-1 is", len(req.Paths[which])-1)
				isReady = true
				//for i := req.Position; i <= len(req.Paths[which])-1; i++ {
				//////////////////////////// the +1 is very important!!!!!!!!!!!!!!!!!!!!!!!
				for i := pos; i <= len(req.Paths[which])-1; i++ {
					//fmt.Println("Inside for. pos is", pos, "which is", which)
					//fmt.Println("Request num", reqNum, "position is", req.Position)
					link := network.GetLinkBetween(req.Paths[which][i], req.Paths[which][i-1])
					//fmt.Println("link is reserved", link.IsReserved)
					if link.IsReserved == false {
						//fmt.Println("link not reserved. Link activation is", link.IsActive)
						//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							//fmt.Println("corresponding reservation.")
							//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
							isReady = isReady && link.IsActive
							cntr++
						}
					}
					if isReady == true {
						//fmt.Println("oppCntr increment. oppCntr is:", oppCntr)
						oppCntr++
					} else {
						break
					}
					if cntr == 0 {
						isReady = false
						break
					}
				}
				if oppCntr >= k {
					//fmt.Println("oppCntr >= k")
					//fmt.Println("Reservation success!")
					//fmt.Println("Run - The req is ", reqNum, "which is ", which)
					whichPath[reqNum] = which
					if changeSrc {
						//fmt.Println("finding position. PositionID is", req.PositionID, "path length is", len(req.Paths[which]))
						req.Position = graph.FindPosition(req.PositionID, req.Paths[which])
					}
					break
				}
			}
		} else {
			//fmt.Println("req.CanMove is true OPP. reqNum is", reqNum)
			//////////// The +1 is very important!!!!!!!!!!!!!!
			for i := req.Position + 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				//fmt.Println("True canMove - Request num", reqNum, "position is", req.Position, "PositionID is", req.PositionID)
				link := network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1])
				//fmt.Println("link is reserved", link.IsReserved)
				if link.IsReserved == false {
					//fmt.Println("link not reserved. Link activation is", link.IsActive)
					//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
					isReady = isReady && link.IsActive
					cntr++
				} else {
					if link.Reservation == reqNum || link.Reservation == -1 {
						//fmt.Println("corresponding reservation.")
						//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					}
				}
				if isReady == true {
					//fmt.Println("oppCntr increment. oppCntr is:", oppCntr)
					oppCntr++
				} else {
					break
				}
				if cntr == 0 {
					isReady = false
					break
				}
			}
		}
		//fmt.Println("checking oppCntr >= k - Request", reqNum, oppCntr >= k)
		//fmt.Println("oppCntr is", oppCntr)
		if oppCntr >= k {
			//fmt.Println("It is Ready!")
			//if !req.CanMove {
			for i := req.Position + 1; i <= req.Position+oppCntr-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
			}
			//}
			req.CanMove = true

			//fmt.Println("Inside noReciveryOPP. reqNum is:", reqNum, "Destination is", req.Dest, "Source is", req.Src, "PositionID is", req.PositionID)
			//for ii, nodee := range req.Paths[whichPath[reqNum]] {
			//	fmt.Println("node index is", ii, "node is", nodee.ID)
			//}

			reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, false)
			numReached += reachedNum
			if w > 0 {
				linksWaiting = append(linksWaiting, w)
			}
		} else if (len(req.Paths[whichPath[reqNum]]) - req.Position - 1) <= oppCntr {
			//fmt.Println("Inside else if!")
			if changeSrc {
				//fmt.Println("finding position. PositionID is", req.PositionID, "path length is", len(req.Paths[whichPath[reqNum]]))
				req.Position = graph.FindPosition(req.PositionID, req.Paths[whichPath[reqNum]])
			}
			//if !req.CanMove {
			for i := req.Position + 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
			}
			//}
			req.CanMove = true

			//fmt.Println("Inside else if no RecoveryOPP. reqNum is:", reqNum, "Destination is", req.Dest, "Source is", req.Src, "PositionID is", req.PositionID)
			//for ii, nodee := range req.Paths[whichPath[reqNum]] {
			//	fmt.Println("node index is", ii, "node is", nodee.ID)
			//}

			reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, false)
			numReached += reachedNum
			if w > 0 {
				linksWaiting = append(linksWaiting, w)
			}
			//fmt.Println("Fill in here. Maybe the remaining links are less than k, but are ready nonetheless.")
		}
		isReady = true
		if changeSrc {
			for m := 1; m <= len(req.Paths[whichPath[reqNum]])-1; m++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).Reservation = -1
			}
			req.CanMove = false
		}
	}
	return numReached, whichPath, linksWaiting
}

func RecoveryRun(network graph.Topology, reqs []*request.Request, whichPath []int, numReached int, runTime int, changeSrc bool) (int, []int, []float64) {
	//numReached := 0
	//fmt.Println("Hiaaa!")
	isReady := true
	willRecover := false
	//whichPath := make([]int, len(reqs))
	linksWaiting := make([]float64, 0)
	var w float64
	var reachedNum int
	var cntr int
	for reqNum, req := range reqs {
		//fmt.Println("Run - The req is: ", reqNum, " The path is: ", whichPath[reqNum])
		if req.HasReached {
			// Release the reserved links
			// Here, req.CanMove is used to release the links, and is set to false
			// to prevent extra work every time the request enters this if statement.
			if req.CanMove {
				//fmt.Println("Req ", reqNum, "Freed the resources.")
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
					//fmt.Println("Freeing link", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
				}
			}
			req.CanMove = false
			continue
		}
		// If the request has no paths, continue until it finds a path.
		if len(req.Paths[0]) == 0 {
			continue
		}
		cntr = 0
		if !req.CanMove {
			//fmt.Println("Can't move!!!!", "req is", reqNum)
			for which, _ := range req.Paths {
				//fmt.Println("Which is ", which, "req is ", reqNum)
				cntr = 0
				isReady = true
				for i := 1; i <= len(req.Paths[which])-1; i++ {
					link := network.GetLinkBetween(req.Paths[which][i], req.Paths[which][i-1])
					if link.IsReserved == false {
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							isReady = isReady && link.IsActive
							cntr++
						} else {
							//fmt.Println("1--- The req is ", reqNum, " It is unfortunately reserved by: ", link.Reservation)
							//fmt.Println("link is", link.ID)
							isReady = false
						}
					}
					if cntr == 0 {
						isReady = false
					}
					if !isReady {
						//recoveryIndex := graph.FindPrecedingRecoveryPoint(req.Paths[which], i, config.GetConfig().GetRecoverySpan())
						///////////////// IMPORTANT!!! CHECK THIS!!!!!
						if !isReady {
							break
						}
					}
					// Solve the isReady issue.
				}
				if isReady {
					whichPath[reqNum] = which
					// changeSrc is for nonoblivious.
					if changeSrc {
						req.Position = graph.FindPosition(req.PositionID, req.Paths[which])
					}
					break
					//fmt.Println("Reservation success!")
					//fmt.Println("Run - The req is ", reqNum, "which is ", which)
				}
			}
		} else {
			if !req.IsRecovering {
				//fmt.Println("Not Recovering!!!!", "req is", reqNum)
				//fmt.Println("req.PositionID is", req.PositionID)
				for i := req.Position + 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					link := network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1])
					if link.IsReserved == false {
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							isReady = isReady && link.IsActive
							cntr++
						} else {
							//fmt.Println("2--- The req is ", reqNum, " It is unfortunately reserved by: ", link.Reservation)
							//fmt.Println("link is", link.ID)
							isReady = false
						}
					}
					if cntr == 0 {
						isReady = false
					}
					//fmt.Println("Before !isReady before willRecover. isReady is", isReady)
					if !isReady {
						//fmt.Println("Checking for recovery!")
						willRecover = checkForRecovery(network, req, reqNum, whichPath[reqNum], req.IsRecovering)
						//fmt.Println("willRecover is", willRecover)
						req.IsRecovering = willRecover
						//isReady = isReady && willRecover
						if willRecover {
							isReady = true
						}
						break
					}
					// Solve the isReady issue.
				}
			} else {
				//fmt.Println("Recovering!!!!", "req is", reqNum)
				isReady = checkForRecovery(network, req, reqNum, whichPath[reqNum], req.IsRecovering)
			}
		}
		//fmt.Println("Request", reqNum, isReady)
		if isReady {
			if !req.IsRecovering {
				// req.CanMove shows the fact that the request has previously reserved the
				// path, and is only trying to swap its way to the end.
				if !req.CanMove {
					//fmt.Println("Req ", reqNum, "is reserving for path ", whichPath[reqNum])
					for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
						network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
						network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
						//fmt.Println("Reserving link: ", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
					}
				}
				req.CanMove = true
				//fmt.Println("------------------------LENGTH IS: ", len(req.Paths[whichPath[reqNum]]), "WHICH IS: ", whichPath[reqNum], "LENGTH OF PATHS IS: ", len(req.Paths))
				reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, false)
				numReached += reachedNum
				if w > 0 {
					linksWaiting = append(linksWaiting, w)
				}
			} else {
				if !req.CanMoveRecovery {
					for i := 1; i <= len(req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex])-1; i++ {
						network.GetLinkBetween(req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i], req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i-1]).IsReserved = true
						network.GetLinkBetween(req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i], req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i-1]).Reservation = reqNum
						//fmt.Println("Reserving link: ", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
					}
				}
				req.CanMoveRecovery = true
				reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, true)
				numReached += reachedNum
				if w > 0 {
					linksWaiting = append(linksWaiting, w)
				}
			}
			if req.HasReached {
				//fmt.Println("Req ", reqNum, " Has reached!")
			}
		}
		isReady = true
		willRecover = false
		// changeSrc is for nonoblivious runs
		if changeSrc {
			for m := 1; m <= len(req.Paths[whichPath[reqNum]])-1; m++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).Reservation = -1
			}
			req.CanMove = false
		}
	}
	return numReached, whichPath, linksWaiting
}

func recoveryRunOPP(network graph.Topology, reqs []*request.Request, whichPath []int, numReached int, runTime int, changeSrc bool) (int, []int, []float64) {
	//fmt.Println("Hiaaa OPP!")
	isReady := true
	willRecover := false
	oppCntr := 0
	k := config.GetConfig().GetOpportunismDegree()
	//whichPath := make([]int, len(reqs))
	linksWaiting := make([]float64, 0)
	var reachedNum int
	var w float64
	var cntr int
	for reqNum, req := range reqs {
		oppCntr = 0
		if req.CanMove {
			for i := 1; i <= req.Position-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
			}
		}
		if req.HasReached {
			// Release the reserved links
			// Here, req.CanMove is used to release the links, and is set to false
			// to prevent extra work every time the request enters this if statement.
			if req.CanMove {
				for i := 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = false
					network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = -1
				}
			}
			req.CanMove = false
			continue
		}
		cntr = 0
		if len(req.Paths[0]) == 0 {
			continue
		}
		// req.Position starts from 1. Check this!!!!!!!!!!!!!!!!!!!!!!!!!
		if !req.CanMove {
			//fmt.Println("Can't move OPP!")
			for which, _ := range req.Paths {
				/////////////////// Fill in here!
				//fmt.Println("Which is ", which, "req is ", reqNum)
				oppCntr = 0
				pos := graph.FindPosition(req.PositionID, req.Paths[which])
				isReady = true
				//for i := req.Position; i <= len(req.Paths[which])-1; i++ {
				//////////////////////////// the +1 is very important!!!!!!!!!!!!!!!!!!!!!!!
				for i := pos; i <= len(req.Paths[which])-1; i++ {
					//fmt.Println("Inside for. pos is", pos, "which is", which)
					//fmt.Println("Request num", reqNum, "position is", req.Position)
					link := network.GetLinkBetween(req.Paths[which][i], req.Paths[which][i-1])
					//fmt.Println("link is reserved", link.IsReserved)
					if link.IsReserved == false {
						//fmt.Println("link not reserved. Link activation is", link.IsActive)
						//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							//fmt.Println("corresponding reservation.")
							//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
							isReady = isReady && link.IsActive
							cntr++
						}
					}
					if isReady == true {
						//fmt.Println("oppCntr increment. oppCntr is:", oppCntr)
						oppCntr++
					} else {
						break
					}
					if cntr == 0 {
						isReady = false
						break
					}
				}
				if oppCntr >= k {
					//fmt.Println("oppCntr >= k")
					//fmt.Println("Reservation success!")
					//fmt.Println("Run - The req is ", reqNum, "which is ", which)
					whichPath[reqNum] = which
					if changeSrc {
						//fmt.Println("finding position. PositionID is", req.PositionID, "path length is", len(req.Paths[which]))
						req.Position = graph.FindPosition(req.PositionID, req.Paths[which])
					}
					break
				}
			}
		} else {
			if !req.IsRecovering {
				//fmt.Println("req.CanMove is true OPP.")
				//fmt.Println("Not Recovering!!!!", "req is", reqNum)
				//fmt.Println("req.PositionID is", req.PositionID)
				//fmt.Println("req.Position is", req.Position)
				//fmt.Println("len(req.Paths[whichPath[reqNum]])-1 is", len(req.Paths[whichPath[reqNum]])-1)
				//////////// The +1 is very important!!!!!!!!!!!!!!
				for i := req.Position + 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
					//fmt.Println("Request num", reqNum, "position is", req.Position)
					link := network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1])
					//fmt.Println("link is reserved", link.IsReserved)
					if link.IsReserved == false {
						//fmt.Println("link not reserved. Link activation is", link.IsActive)
						//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
						isReady = isReady && link.IsActive
						cntr++
					} else {
						if link.Reservation == reqNum || link.Reservation == -1 {
							//fmt.Println("corresponding reservation.")
							//fmt.Println("Link is", link.ID, "link.IsActive is", link.IsActive)
							isReady = isReady && link.IsActive
							cntr++
						}
					}
					if !isReady {
						//fmt.Println("Checking for recovery!")
						willRecover = checkForRecovery(network, req, reqNum, whichPath[reqNum], req.IsRecovering)
						//fmt.Println("willRecover is", willRecover)
						req.IsRecovering = willRecover
						//isReady = isReady && willRecover
						if willRecover {
							isReady = true
						}
						break
					}
					if isReady == true {
						//fmt.Println("oppCntr increment. oppCntr is:", oppCntr)
						oppCntr++
					} else {
						break
					}
					if cntr == 0 {
						isReady = false
						break
					}
				}
			} else {
				//fmt.Println("Recovering!!!!", "req is", reqNum)
				isReady = checkForRecovery(network, req, reqNum, whichPath[reqNum], req.IsRecovering)
			}
		}
		//fmt.Println("Request", reqNum, oppCntr >= k)
		if oppCntr >= k {
			//fmt.Println("It is Ready!")
			//if !req.CanMove {
			for i := req.Position + 1; i <= req.Position+oppCntr-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
			}
			//}
			req.CanMove = true

			//fmt.Println("Inside noReciveryOPP. reqNum is:", reqNum, "Destination is", req.Dest, "Source is", req.Src, "PositionID is", req.PositionID)
			//for ii, nodee := range req.Paths[whichPath[reqNum]] {
			//	fmt.Println("node index is", ii, "node is", nodee.ID)
			//}

			reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, false)
			numReached += reachedNum
			if w > 0 {
				linksWaiting = append(linksWaiting, w)
			}
		} else if (len(req.Paths[whichPath[reqNum]]) - req.Position - 1) <= oppCntr {
			if changeSrc {
				//fmt.Println("finding position. PositionID is", req.PositionID, "path length is", len(req.Paths[whichPath[reqNum]]))
				req.Position = graph.FindPosition(req.PositionID, req.Paths[whichPath[reqNum]])
			}
			//if !req.CanMove {
			for i := req.Position + 1; i <= len(req.Paths[whichPath[reqNum]])-1; i++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).IsReserved = true
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).Reservation = reqNum
			}
			//}
			req.CanMove = true

			//fmt.Println("Inside else if no RecoveryOPP. reqNum is:", reqNum, "Destination is", req.Dest, "Source is", req.Src, "PositionID is", req.PositionID)
			//for ii, nodee := range req.Paths[whichPath[reqNum]] {
			//	fmt.Println("node index is", ii, "node is", nodee.ID)
			//}

			reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, false)
			numReached += reachedNum
			if w > 0 {
				linksWaiting = append(linksWaiting, w)
			}
			//fmt.Println("Fill in here. Maybe the remaining links are less than k, but are ready nonetheless.")
		} else {
			//fmt.Println("recoveryRunOpp - Before reserving. reqNum is", reqNum)
			//fmt.Println("recoveryRunOpp - Before reserving. req.Src.ID is", req.Src.ID)
			//fmt.Println("recoveryRunOpp - Before reserving. req.Dest.ID is", req.Dest.ID)
			//fmt.Println("recoveryRunOpp - Before reserving. req.PositionID is", req.PositionID)
			if req.IsRecovering {
				//fmt.Println("req.IsRecovering is true.")
				if isReady {
					//fmt.Println("isReady is true.")
					if !req.CanMoveRecovery {
						for i := 1; i <= len(req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex])-1; i++ {
							network.GetLinkBetween(req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i], req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i-1]).IsReserved = true
							network.GetLinkBetween(req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i], req.RecoveryPaths[whichPath[reqNum]][req.RecoveryPathCursor][req.RecoveryPathIndex][i-1]).Reservation = reqNum
							//fmt.Println("Reserving link: ", network.GetLinkBetween(req.Paths[whichPath[reqNum]][i], req.Paths[whichPath[reqNum]][i-1]).ID, "WhichPath is ", whichPath[reqNum])
						}
					}
					req.CanMoveRecovery = true
					//fmt.Println("recoveryRunOpp - Before quantum.ES()")
					reachedNum, w = quantum.ES(req, network, runTime, whichPath[reqNum], changeSrc, true)
					numReached += reachedNum
					if w > 0 {
						linksWaiting = append(linksWaiting, w)
					}
					//fmt.Println("recoveryRunOpp - after quantum.ES(). reqNum is", reqNum)
					//fmt.Println("recoveryRunOpp - after quantum.ES(). req.Position is", req.Position)
					//fmt.Println("recoveryRunOpp - after quantum.ES(). req.PositionID is", req.PositionID)
				}
			}
		}
		isReady = true
		willRecover = false
		if changeSrc {
			for m := 1; m <= len(req.Paths[whichPath[reqNum]])-1; m++ {
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).IsReserved = false
				network.GetLinkBetween(req.Paths[whichPath[reqNum]][m], req.Paths[whichPath[reqNum]][m-1]).Reservation = -1
			}
			req.CanMove = false
		}
	}
	return numReached, whichPath, linksWaiting
}

func checkForRecovery(network graph.Topology, req *request.Request, reqNum int, pathIndex int, isRecovering bool) bool {
	isRecPoint := false
	willRecover := true
	cntr := 0
	if !isRecovering {
		for _, recIndex := range graph.NumRecoveryIndex(len(req.Paths[pathIndex])) {
			if req.Position-1 == recIndex {
				isRecPoint = true
				break
			}
		}
		if !isRecPoint {
			return false
		}
		if len(req.RecoveryPaths[pathIndex]) == 0 {
			return false
		}
		//fmt.Println("length of req.RecoveryPaths is", len(req.RecoveryPaths))
		//fmt.Println("length of req.RecoveryPaths[pathIndex] is", len(req.RecoveryPaths[pathIndex]))
		//fmt.Println("pathIndex is", pathIndex)
		//fmt.Println("req.Src.ID is", req.Src.ID, "req.Dest.ID is", req.Dest.ID)
		//fmt.Println("req.PositionID is", req.PositionID)
		cursors := graph.NumRecoveryIndex(len(req.Paths[pathIndex]))
		var cursorIndex int
		for index, val := range cursors {
			if val == req.Position-1 {
				cursorIndex = index
			}
		}
		//fmt.Println("cursorIndex is", cursorIndex)
		for recPathIndex, recPath := range req.RecoveryPaths[pathIndex][cursorIndex] {
			cntr = 0
			willRecover = true
			if len(recPath) == 0 {
				continue
			}
			for i := 1; i <= len(recPath)-1; i++ {
				link := network.GetLinkBetween(recPath[i], recPath[i-1])
				if link.IsReserved == false {
					willRecover = willRecover && link.IsActive
					cntr++
				} else {
					if link.Reservation == reqNum || link.Reservation == -1 {
						willRecover = willRecover && link.IsActive
						cntr++
					} else {
						willRecover = false
					}
				}
				if cntr == 0 {
					willRecover = false
				}
				if !willRecover {
					break
				}
			}
			if cntr == 0 {
				willRecover = false
			}
			if willRecover {
				//whichRec = recPathIndex
				req.RecoveryPathIndex = recPathIndex
				req.RecoveryPathCursor = cursorIndex
				return true
			}
		}
	} else {
		for i := req.RecoveryPosition + 1; i <= len(req.RecoveryPaths[pathIndex][req.RecoveryPathCursor][req.RecoveryPathIndex])-1; i++ {
			link := network.GetLinkBetween(req.RecoveryPaths[pathIndex][req.RecoveryPathCursor][req.RecoveryPathIndex][i], req.RecoveryPaths[pathIndex][req.RecoveryPathCursor][req.RecoveryPathIndex][i-1])
			if link.IsReserved == false {
				willRecover = willRecover && link.IsActive
				cntr++
			} else {
				if link.Reservation == reqNum || link.Reservation == -1 {
					willRecover = willRecover && link.IsActive
					cntr++
				} else {
					willRecover = false
				}
			}
			if cntr == 0 {
				willRecover = false
			}
			if !willRecover {
				break
			}
		}
		return willRecover
	}
	return false
}

func findRecoveryPaths(reqs []*request.Request, network graph.Topology) {
	//recoverySpan := config.GetConfig().GetRecoverySpan()
	recoveryHasContention := config.GetConfig().GetRecoveryHasContention()
	//mainSegment := make([]*graph.Link, recoverySpan)
	if !recoveryHasContention {
		for _, req := range reqs {
			for p, _ := range req.Paths {
				linksToPrune := path.PathToLinks(req.Paths[p], network)
				graph.Prune(linksToPrune)
			}
		}
	}
	recAgg := config.GetConfig().GetRecoveryAggressiveness()
	recSpan := config.GetConfig().GetRecoverySpan()
	for _, req := range reqs {
		pathPass := false
		for pathIndex, p := range req.Paths {
			if !pathPass {
				req.RecoveryPaths[pathIndex] = make([][][]*graph.Node, 1)
				pathPass = true
			} else {
				dummy := make([][][]*graph.Node, 1)
				req.RecoveryPaths = append(req.RecoveryPaths, dummy)
			}
			isInitiated := false
			for cursor, recoveryIndex := range graph.NumRecoveryIndex(len(p)) {
				if isInitiated {
					temp := make([][]*graph.Node, 1)
					req.RecoveryPaths[pathIndex] = append(req.RecoveryPaths[pathIndex], temp)
				}
				if cursor == len(graph.NumRecoveryIndex(len(p)))-1 {
					break
				}
				var linksToPrune []*graph.Link
				//fmt.Println("recoveryIndex is:", recoveryIndex)
				//fmt.Println("recSpan is:", recSpan)
				if recoveryIndex == 0 {
					if recoveryIndex+recSpan+1 <= len(p) {
						linksToPrune = path.PathToLinks(p[recoveryIndex:recoveryIndex+recSpan+1], network)
					} else {
						linksToPrune = path.PathToLinks(p[recoveryIndex:len(p)], network)
					}
				} else {
					if recoveryIndex+recSpan+1 <= len(p) {
						linksToPrune = path.PathToLinks(p[recoveryIndex-1:recoveryIndex+recSpan+1], network)
					} else {
						linksToPrune = path.PathToLinks(p[recoveryIndex-1:len(p)], network)
					}
				}
				// At the end of path.PF, the whole network is depruned.
				graph.Prune(linksToPrune)
				//log.PrintLinks(linksToPrune)
				auxiliaryReq := new(request.Request)
				request.CopyRequest(auxiliaryReq, req)
				//auxiliaryReq := req
				auxiliaryReq.Src = p[recoveryIndex]
				if recoveryIndex+recSpan < len(p) {
					auxiliaryReq.Dest = p[recoveryIndex+recSpan]
				} else {
					auxiliaryReq.Dest = p[len(p)-1]
				}
				//fmt.Println("auxiliaryReq.Src is:", auxiliaryReq.Src.ID)
				//fmt.Println("auxiliaryReq.Dest is:", auxiliaryReq.Dest.ID)
				//conf := config.GetConfig()
				//confPointer := &conf
				//confPointer.SetAggressiveness(recAgg)
				config.SetAggressiveness(recAgg)
				//fmt.Println("aggressiveness is:", config.GetConfig().GetAggressiveness())
				requests := make([]*request.Request, 1)
				requests[0] = auxiliaryReq
				path.PF(network, requests, "modified greedy", true)
				cntr := 0
				//if len(req.RecoveryPaths[pathIndex] == 1) {
				//	req.RecoveryPaths[pathIndex][cursor] := make([][]*graph.Node, 1)
				//}
				req.RecoveryPaths[pathIndex][cursor] = make([][]*graph.Node, 1)
				//fmt.Println("First - length of req.RecoveryPaths[pathIndex][cursor] is:", len(req.RecoveryPaths[pathIndex][cursor]))
				hasPassed := false
				for i := 0; i < int(math.Min(float64(recAgg), float64(len(requests[0].Paths)))); i++ {
					if len(requests[0].Paths[i]) == 0 {
						//fmt.Println("Hi! recovery path is nil.")
						continue
					}
					//fmt.Println("length of req.RecoveryPaths[pathIndex][cursor] is:", len(req.RecoveryPaths[pathIndex][cursor]))
					if !hasPassed {
						//fmt.Println("pathIndex is:", pathIndex, "cursor is:", cursor, "cntr is:", cntr)
						req.RecoveryPaths[pathIndex][cursor][cntr] = requests[0].Paths[i]
					} else {
						req.RecoveryPaths[pathIndex][cursor] = append(req.RecoveryPaths[pathIndex][cursor], requests[0].Paths[i])
					}
					//if len(req.RecoveryPaths[pathIndex][cursor]) == 1 {
					//	hasPassed = false
					//} else {
					hasPassed = true
					//}
					//req.RecoveryPaths[pathIndex][cursor][cntr] = requests[0].Paths[i]
					cntr++
				}
				//if hasPassed {
				isInitiated = true
				//} else {

				//}
			}
		}
	}
}

// Each profile will have a unique profile id.
//func BuildProfile(profileID int) profile {

//}
