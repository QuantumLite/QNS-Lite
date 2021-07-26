package benchmark

import (
	"fmt"

	"example.com/profile"
	"example.com/request"
)

type Benchmarker struct {
	keepReqs       bool
	regenerateReqs bool
	refreshSources bool
	// ignoreLeftOvers deals with the leftover requests when looking for paths. It allows us to
	// prevent infinite loops.
	ignoreLeftOvers     bool
	Throughput          []float64
	TotalWaitingTime    []int
	LinksWaitingTime    []float64
	reqsWaitingTime     [][]int
	priorityWaitingTime [][]int
	profile             profile.Profile
	reqs                []*request.Request
}

func (bm *Benchmarker) Set(itr int, prof string, topology string) {
	bm.Throughput = make([]float64, itr)
	bm.TotalWaitingTime = make([]int, itr)
	bm.LinksWaitingTime = make([]float64, 0)
	if prof == profile.MODIFIED_GREEDY {
		mgp := new(profile.ModifiedGreedyProfile)
		mgp.Build(topology)
		bm.profile = mgp
		bm.keepReqs = false
		bm.regenerateReqs = false
		bm.refreshSources = false
		bm.ignoreLeftOvers = false
	} else if prof == profile.NONOBLIVIOUS_LOCAL {
		nol := new(profile.NonObliviousLocal)
		nol.Build(topology)
		bm.profile = nol
		bm.keepReqs = false
		bm.regenerateReqs = false
		bm.refreshSources = true
		bm.ignoreLeftOvers = true
	} else if prof == profile.QPASS {
		qpass := new(profile.QPass)
		qpass.Build(topology)
		bm.profile = qpass
		bm.keepReqs = false
		bm.regenerateReqs = false
		bm.refreshSources = false
		bm.ignoreLeftOvers = false
	} else {
		fmt.Println("Benchmark: Caution! The profile is not implemented.")
	}
	bm.priorityWaitingTime = make([][]int, bm.profile.GetPriorityLen())
	for i := 0; i < bm.profile.GetPriorityLen(); i++ {
		bm.priorityWaitingTime[i] = make([]int, 0)
	}
}

func (bm *Benchmarker) Start(itr int, maxItr int) {
	///////////////////////// This might be unnecessary, since now we have the regenerateReqs()
	///////////////////////// function.
	bm.priorityWaitingTime = make([][]int, bm.profile.GetPriorityLen())
	for i := 0; i < bm.profile.GetPriorityLen(); i++ {
		bm.priorityWaitingTime[i] = make([]int, 0)
	}
	if !bm.keepReqs {
		//reqs := profile.GenRequests(config.GetConfig().GetNumRequests(), bm.profile.GetNetwork(), config.GetConfig().GetIsMultiPath(), bm.profile.GetPathAlgorithm(), bm.ignoreLeftOvers)
		reqs := bm.profile.GenRequests(bm.ignoreLeftOvers)
		bm.reqs = reqs
	}
	for i := 0; i <= itr-1; i++ {
		//fmt.Println(*bm)
		//fmt.Println("Iteration", i)
		//if bm.refreshSources {
		//	for m, _ := range bm.reqs {
		//		bm.reqs[m].Src = bm.reqs[m].InitialSrc
		//	}
		//}
		bm.profile.Run(bm.reqs, maxItr)
		//fmt.Println(*bm)
		bm.TotalWaitingTime[i] = bm.profile.GetRunTime()
		bm.LinksWaitingTime = bm.profile.GetLinksWaitingTime()
		for reqIndex, Req := range bm.reqs {
			bm.reqsWaitingTime[reqIndex][i] = Req.ServingTime - Req.GenerationTime
			for i := 1; i <= bm.profile.GetPriorityLen(); i++ {
				if Req.Priority == i {
					bm.priorityWaitingTime[i-1] = append(bm.priorityWaitingTime[i-1], int((Req.ServingTime-Req.GenerationTime)/(Req.Priority*bm.profile.GetNetwork().Distance(Req.Src, Req.Dest, "hop"))))
					break
				}
			}
		}
		bm.profile.Clear()
		//bm.LinksWaitingTime = make([]float64, 0)
		for _, req := range bm.reqs {
			request.ClearReq(req)
			if bm.refreshSources {
				req.Src = req.InitialSrc
			}
		}
		//fmt.Println("GOTCHA!", *bm)
	}
	bm.profile.Clear()
	bm.profile.Stop()
}

func (bm *Benchmarker) SetKeepReqs(keepReqs bool) {
	bm.keepReqs = keepReqs
}

func (bm *Benchmarker) RegenerateReqs(itr int) {
	//bm.reqs = bm.profile.GenRequests(config.GetConfig().GetNumRequests(), bm.profile.GetNetwork(), config.GetConfig().GetIsMultiPath(), bm.profile.GetPathAlgorithm(), bm.ignoreLeftOvers)
	bm.reqs = bm.profile.GenRequests(bm.ignoreLeftOvers)
	bm.reqsWaitingTime = make([][]int, len(bm.reqs))
	for i := 0; i < len(bm.reqs); i++ {
		bm.reqsWaitingTime[i] = make([]int, itr)
	}
}

func (bm *Benchmarker) AverageWaiting(maxItr int) float64 {
	sum := 0
	meanLength := len(bm.TotalWaitingTime)
	for _, val := range bm.TotalWaitingTime {
		if val >= maxItr-1 {
			meanLength -= 1
			continue
		}
		sum += val
	}
	return float64(sum) / float64(meanLength)
}

func (bm *Benchmarker) VarianceWaiting(maxItr int) float64 {
	sum := float64(0)
	ave := bm.AverageWaiting(maxItr)
	varLength := len(bm.TotalWaitingTime)
	for _, val := range bm.TotalWaitingTime {
		if val >= maxItr-1 {
			varLength--
			continue
		}
		sum += (float64(val) - ave) * (float64(val) - ave)
	}
	return float64(sum) / float64(varLength)
}

func (bm *Benchmarker) PriorityAverageWaiting(maxItr int) []float64 {
	sum := 0
	means := make([]float64, bm.profile.GetPriorityLen())
	meanLength := make([]int, bm.profile.GetPriorityLen())
	for i := 0; i < bm.profile.GetPriorityLen(); i++ {
		sum = 0
		meanLength[i] = len(bm.priorityWaitingTime[i])
		for _, val := range bm.priorityWaitingTime[i] {
			if val >= maxItr-1 {
				meanLength[i] -= 1
				continue
			}
			sum += val
		}
		means[i] = float64(sum) / float64(meanLength[i])
	}

	return means
}
