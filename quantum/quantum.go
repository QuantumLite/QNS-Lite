package quantum

import (
	"fmt"
	"math/rand"
	"time"

	"example.com/config"
	"example.com/graph"
	"example.com/path"
	"example.com/request"
)

type swapper interface {
}

type generator interface {
	genEntanglement(links []*graph.Link)
}

var p_gen, p_swap float64

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	p_gen = config.GetConfig().GetPGen()
	fmt.Println("quantum init. p_gen is", p_gen)
	p_swap = config.GetConfig().GetPSwap()
}

func genEntanglement(links []*graph.Link, roundNum int) {
	lifetime := config.GetConfig().GetLifetime()
	p_gen = config.GetConfig().GetPGen()
	//var r float64
	for _, link := range links {
		if link.IsActive == true {
			if link.Age == lifetime {
				link.IsActive = false
			} else {
				link.Age++
			}
			continue
		}
		if link.IsActive == false {
			//r = rand.Float64()
			if rand.Float64() <= p_gen {
				//fmt.Println("Link successfully generated!")
				link.IsActive = true
				link.Age = 0
				link.GenerationTime = roundNum
			}
		}
	}
}

func genEntangledLink(link *graph.Link) {
	if rand.Float64() <= p_gen {
		link.IsActive = true
	}
}

// TODO: Write a less dumb swap!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func swap(req *request.Request, path path.Path, network graph.Topology, roundNum int, changeSrc bool) (bool, float64) {
	linkWaiting := float64(-1)
	//fmt.Println(req.Position)
	if req.Position == len(path)-1 {
		//fmt.Println("SWAP: HAS REACHED 1")
		req.HasReached = true
		req.ServingTime = roundNum
		return true, linkWaiting
	}
	if graph.IsEqual(req.Dest.ID, req.PositionID) {
		//fmt.Println("Inside quantum.swap, has reached: Request destination is:", req.Dest.ID, "request position is:", req.PositionID, "Last node of path is:", path[len(path)-1].ID)
		req.HasReached = true
		req.ServingTime = roundNum
		return true, linkWaiting
	}
	///////////////////// If swapping is unsuccessful, do the links get destroyed?
	///////////////////// If swapping is unsuccessful, shall we unreserve the links, as
	///////////////////// we have done here?
	if len(path) == 1 {
		//fmt.Println("Dumb path with single node. req is", req.Src.ID, req.Dest.ID, "the single node is", path[0].ID)
		if graph.IsEqual(req.Dest.ID, path[0].ID) {
			//fmt.Println("Dumb path has reached its destination.")
			//fmt.Println("Dumb path req.PositionID is", req.PositionID)
			req.HasReached = true
			req.ServingTime = roundNum
			return true, linkWaiting
		}
	}
	if req.Position == 1 {
		if len(path) == 1 {
			//fmt.Println("Inside quantum.swap: Request destination is:", req.Dest.ID, "request position is:", req.PositionID, "Last node of path is:", path[len(path)-1].ID)
		}
		if len(path) > 1 {
			network.GetLinkBetween(path[req.Position-1], path[req.Position]).IsActive = false
		}
	}
	//fmt.Println("Inside quantum swap. Position is:", req.Position, "PositionID is", req.PositionID, "path length is:", len(path))
	if len(path) > 1 {
		network.GetLinkBetween(path[req.Position], path[req.Position+1]).IsActive = false
	}
	if rand.Float64() <= p_swap {
		//fmt.Println("SWAP: Swap successful.")
		if req.Position == 1 {
			network.GetLinkBetween(path[req.Position-1], path[req.Position]).IsReserved = false
			network.GetLinkBetween(path[req.Position-1], path[req.Position]).Reservation = -1
			network.GetLinkBetween(path[req.Position-1], path[req.Position]).ConsumptionTime = roundNum
		}
		network.GetLinkBetween(path[req.Position], path[req.Position+1]).IsReserved = false
		network.GetLinkBetween(path[req.Position], path[req.Position+1]).Reservation = -1
		network.GetLinkBetween(path[req.Position], path[req.Position+1]).ConsumptionTime = roundNum
		linkWaiting = float64(network.GetLinkBetween(path[req.Position-1], path[req.Position]).ConsumptionTime - network.GetLinkBetween(path[req.Position-1], path[req.Position]).GenerationTime)
		//if req.Position == 1 {
		//	network.GetLinkBetween(path[req.Position-1], path[req.Position]).IsActive = false
		//}
		//network.GetLinkBetween(path[req.Position], path[req.Position+1]).IsActive = false
		//req.Position++
		req.Position++
		////////////////////////////////// VERY IMPORTANT!!!!! THE -1 IS NECESSARY HERE!!!
		req.PositionID = path[req.Position-1].ID
		if changeSrc {
			////////////////////////////////////// VERY IMPORTANT!!!!! RELATED TO THE
			////////////////////////////////////// VERY IMPORTANT COMMENT ABOVE!
			req.Src = path[req.Position-1]
		}
		if req.Position == len(path)-1 {
			//fmt.Println("SWAP: HAS REACHED 2")
			req.HasReached = true
			req.ServingTime = roundNum
			return true, linkWaiting
		}
		//fmt.Println("Quantum: Dest is:", req.Dest.ID, "PositionID is:", req.PositionID)
		if graph.IsEqual(req.Dest.ID, req.PositionID) {
			//fmt.Println("Inside quantum.swap, has reached: Request destination is:", req.Dest.ID, "request position is:", req.PositionID, "Last node of path is:", path[len(path)-1].ID)
			req.HasReached = true
			req.ServingTime = roundNum
			return true, linkWaiting
		}
	} else {
		/////////// Very important! Check this! dramatically decreases performance!
		req.Position = 1
		req.Src = req.InitialSrc
		//fmt.Println("swap failed!")
	}
	return false, linkWaiting
}

func EG(links []*graph.Link, roundNum int) {
	genEntanglement(links, roundNum)
}

func ES(req *request.Request, network graph.Topology, roundNum int, whichPath int, changeSrc bool, isRecovery bool) (int, float64) {
	var reached bool
	var linkWaiting float64
	numReached := 0
	if !isRecovery {
		reached, linkWaiting = swap(req, req.Paths[whichPath], network, roundNum, changeSrc)
	} else {
		auxiliaryReq := new(request.Request)
		request.CopyRequest(auxiliaryReq, req)
		auxiliaryReq.Position = req.RecoveryPosition
		//fmt.Println("quantum ES - Before recovery swap. req.PositionID is", req.PositionID)
		//fmt.Println("quantum ES - Before recovery swap. whichPath is", whichPath, "req.RecoveryPathCursor is", req.RecoveryPathCursor, "req.RecoveryPathIndex is", req.RecoveryPathIndex)
		//for _, nodede := range req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex] {
		//	fmt.Println("quantum ES - Before recovery swap - The recovery path is", nodede.ID)
		//}
		//fmt.Println("quantum ES - Before recovery swap. auxiliaryReq.Position is", auxiliaryReq.Position)
		reached, linkWaiting = swap(auxiliaryReq, req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex], network, roundNum, changeSrc)
		//fmt.Println("quantum ES - after recovery swap. auxiliaryReq.Position is", auxiliaryReq.Position)
		//fmt.Println("quantum ES - after recovery swap. auxiliaryReq.PositionID is", auxiliaryReq.PositionID)
		//fmt.Println("quantum ES - after recovery swap. req.Dest.ID is", req.Dest.ID)
		//if !graph.IsEqual(auxiliaryReq.PositionID, req.Dest.ID) {
		if !graph.IsEqual(auxiliaryReq.PositionID, req.Paths[whichPath][len(req.Paths[whichPath])-1].ID) {
			reached = false
			req.HasReached = false
		} else {
			reached = true
			req.HasReached = true
		}
		//req.HasReached = auxiliaryReq.HasReached
		req.PositionID = auxiliaryReq.PositionID
		req.RecoveryPosition = auxiliaryReq.Position
		if auxiliaryReq.Position == len(req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex])-1 {
			///////////// The next line is added to compensate for the positionID when the
			///////////// recovery is finished successfully, since the current system
			///////////// falls one index behind when handling position and positionID.
			//fmt.Println("AHAAAAAAAAAAAAAAAAAA - last recovery finished.")
			req.PositionID = req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex][len(req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex])-1].ID
			req.IsRecovering = false
			req.CanMoveRecovery = false
			//fmt.Println("quantum ES - checking the auxiliaryReq. auxiliaryReq.Position is", auxiliaryReq.Position)
			//fmt.Println("quantum ES - Checking the auxiliaryReq. req.PositionID is", req.PositionID)
			//fmt.Println("quantum ES - Checking the auxiliaryReq. req.PositionID is", req.PositionID)
			//for _, nodede := range req.Paths[whichPath] {
			//	fmt.Println("The path is", nodede.ID)
			//}
			//for _, nodede := range req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex] {
			//	fmt.Println("The recovery path is", nodede.ID)
			//}
			req.Position = graph.FindPosition(req.PositionID, req.Paths[whichPath])
			req.RecoveryPosition = 1
			//fmt.Println("quantum ES - req.Position after graph.FindPosition is", req.Position)
			graph.UnreservePath(req.RecoveryPaths[whichPath][req.RecoveryPathCursor][req.RecoveryPathIndex], network)
		}
		if req.Position == len(req.Paths[whichPath])-1 {
			req.HasReached = true
			req.ServingTime = roundNum
			reached = true
		}
	}
	if reached == true {
		numReached++
	}
	return numReached, linkWaiting
}
