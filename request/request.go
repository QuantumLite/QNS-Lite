package request

import (
	"errors"
	"fmt"
	"math/rand"

	"example.com/config"
	"example.com/graph"
)

type Request struct {
	Src                *graph.Node
	InitialSrc         *graph.Node
	Dest               *graph.Node
	Paths              [][]*graph.Node
	RecoveryPaths      [][][][]*graph.Node
	PositionID         []int
	Position           int
	RecoveryPosition   int
	RecoveryPathIndex  int
	RecoveryPathCursor int
	Priority           int
	GenerationTime     int
	ServingTime        int
	HasReached         bool
	CanMove            bool
	CanMoveRecovery    bool
	IsRecovering       bool
}

/////////////////////////// This package should be more general. Use an interface to avoid
/////////////////////////// many if-else blocks.

///////////////////////////////////////////// Change this function to receive requests and append to them.
func genRequests(N int, ids [][]int, priority []int, topology string, roundNum int) ([]*Request, error) {
	if topology != graph.GRID && topology != graph.RING {
		return nil, errors.New("request.genRequests: The requested topology is not valid!")
	}
	var reqs []*Request
	if topology == graph.GRID {
		var isSame bool
		var r [2]int
		reqs = make([]*Request, N)
		//s := make([]int, len(ids[0]))
		//d := make([]int, len(ids[0]))
		for i := 0; i < N; i++ {
			isSame = true
			for isSame == true {
				r[0] = rand.Intn(len(ids))
				r[1] = rand.Intn(len(ids))
				isSame = r[0] == r[1]
			}
			reqs[i] = new(Request)
			reqs[i].Src = new(graph.Node)
			reqs[i].InitialSrc = new(graph.Node)
			reqs[i].Dest = new(graph.Node)
			reqs[i].Src.ID = make([]int, 2)
			reqs[i].InitialSrc.ID = make([]int, 2)
			reqs[i].Dest.ID = make([]int, 2)
			reqs[i].Src = graph.MakeNode(ids[r[0]], config.GetConfig().GetMemory())
			reqs[i].InitialSrc = graph.MakeNode(ids[r[0]], config.GetConfig().GetMemory())
			reqs[i].Dest = graph.MakeNode(ids[r[1]], config.GetConfig().GetMemory())
			reqs[i].PositionID = make([]int, 2)
			reqs[i].Position = 1
			reqs[i].RecoveryPosition = 1
			reqs[i].RecoveryPathIndex = 0
			reqs[i].RecoveryPathCursor = 0
			reqs[i].Priority = priority[i]
			reqs[i].GenerationTime = roundNum
			reqs[i].HasReached = false
			reqs[i].CanMove = false
			reqs[i].CanMoveRecovery = false
			reqs[i].IsRecovering = false

			////////////// TODO: IMPORTANT!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			reqs[i].Paths = make([][]*graph.Node, 1)
			reqs[i].RecoveryPaths = make([][][][]*graph.Node, 1)
		}
	} else if topology == graph.RING {
		var isSame bool
		var r [2]int
		reqs = make([]*Request, N)
		for i := 0; i < N; i++ {
			isSame = true
			for isSame == true {
				r[0] = rand.Intn(len(ids))
				r[1] = rand.Intn(len(ids))
				isSame = r[0] == r[1]
			}
			reqs[i] = new(Request)
			reqs[i].Src = new(graph.Node)
			reqs[i].InitialSrc = new(graph.Node)
			reqs[i].Dest = new(graph.Node)
			reqs[i].Src.ID = make([]int, 1)
			reqs[i].InitialSrc.ID = make([]int, 1)
			reqs[i].Dest.ID = make([]int, 1)
			reqs[i].Src = graph.MakeNode(ids[r[0]], config.GetConfig().GetMemory())
			reqs[i].InitialSrc = graph.MakeNode(ids[r[0]], config.GetConfig().GetMemory())
			reqs[i].Dest = graph.MakeNode(ids[r[1]], config.GetConfig().GetMemory())
			reqs[i].PositionID = make([]int, 1)
			reqs[i].Position = 1
			reqs[i].RecoveryPosition = 1
			reqs[i].RecoveryPathIndex = 0
			reqs[i].RecoveryPathCursor = 0
			reqs[i].Priority = priority[i]
			reqs[i].GenerationTime = roundNum
			reqs[i].HasReached = false
			reqs[i].CanMove = false
			reqs[i].CanMoveRecovery = false
			reqs[i].IsRecovering = false

			////////////// TODO: IMPORTANT!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			reqs[i].Paths = make([][]*graph.Node, 1)
			reqs[i].RecoveryPaths = make([][][][]*graph.Node, 1)
		}
	} else {
		fmt.Println("Request - genRequests: Caution! Input topology is not supported.")
		return nil, nil
	}
	return reqs, nil
}

// TODO: Handle the priorities
func RG(N int, ids [][]int, priority []int, topology string, roundNum int) ([]*Request, error) {
	return genRequests(N, ids, priority, topology, roundNum)
}

//func isPathless(req *Request) bool {
//	if len(req.Paths) > 1 {
//		return false
//	} else {
//		if Paths[0][0].Memory = 0 {
//			return true
//		}
//	}
//	return false
//}

func ClearReq(req *Request) {
	req.Position = 1
	req.PositionID = req.Src.ID
	req.RecoveryPosition = 1
	req.RecoveryPathIndex = 0
	req.RecoveryPathCursor = 0
	req.HasReached = false
	req.CanMove = false
	req.CanMoveRecovery = false
	req.IsRecovering = false
}

func ClearReqPaths(reqs []*Request) {
	for _, req := range reqs {
		req.Paths = make([][]*graph.Node, 1)
		/// Is this necessary?
		req.RecoveryPaths = make([][][][]*graph.Node, 1)
	}
}

func CopyRequest(newReq, reqToCopy *Request) {
	newReq.Src = reqToCopy.Src
	newReq.InitialSrc = reqToCopy.InitialSrc
	newReq.Dest = reqToCopy.Dest
	newReq.PositionID = reqToCopy.PositionID
	newReq.Position = reqToCopy.Position
	newReq.RecoveryPosition = reqToCopy.RecoveryPosition
	newReq.RecoveryPathIndex = reqToCopy.RecoveryPathIndex
	newReq.RecoveryPathCursor = reqToCopy.RecoveryPathCursor
	newReq.Priority = reqToCopy.Priority
	newReq.GenerationTime = reqToCopy.GenerationTime
	newReq.HasReached = reqToCopy.HasReached
	newReq.CanMove = reqToCopy.CanMove
	newReq.CanMoveRecovery = reqToCopy.CanMoveRecovery
	newReq.IsRecovering = reqToCopy.IsRecovering
	newReq.Paths = make([][]*graph.Node, 1)
	newReq.RecoveryPaths = make([][][][]*graph.Node, 1)
}

// GatherRemainingRequests() gathers the requests not
func GatherRemainingRequests() {

}
