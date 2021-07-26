package main

import (
	"fmt"

	"example.com/benchmark"
	"example.com/config"
)

func main() {
	//var grid *graph.Grid = new(graph.Grid)
	//grid.Build()
	//ids := grid.GetNodeIDs()
	//fmt.Println("IDs", ids)
	//links := grid.GetLinks()
	//fmt.Println("Links", links)
	//fmt.Println("Links", links[2][0][1].ID)
	//num := 5
	//var priority []int
	//priority = make([]int, num)
	//for i := 0; i < num; i++ {
	//	priority[i] = 1
	//}
	//reqs, err := request.RG(num, ids, priority, "grid")
	//for i := 0; i < num; i++ {
	//	fmt.Println(i, reqs[i].Src)
	//	fmt.Println(reqs[i].Dest)
	//	for _, node := range reqs[i].Paths[0] {
	//		fmt.Println("PATHS FOR THIS REQUEST", node.ID)
	//	}
	//}
	//fmt.Println(reqs)
	//fmt.Println("Request generation error:", err)
	//path.PF(grid, reqs, "modified greedy", false)
	//check := 2
	//fmt.Println(reqs[check].Src, reqs[check].Dest)
	//for _, link := range links {
	//	fmt.Println("Before", link.IsActive)
	//}
	//quantum.EG(links)
	//for _, link := range links {
	//	fmt.Println("After", link.IsActive)
	//}

	/////////////////////////////////////////////////////// Implement lifetime!!!!!!!!

	/*
		itr := 100
		maxItr := 5000
		bm := new(benchmark.Benchmarker)
		bm.Set(itr, "modified greedy", "grid")
		bm.Start(itr, maxItr)
		fmt.Println(*bm)
		fmt.Println("The average waiting time is:", bm.AverageWaiting(maxItr))
		fmt.Println("The variance of the waiting time is:", bm.VarianceWaiting(maxItr))
		config.SetOpportunism(true)
		bm.SetKeepReqs(true)
		bm.Start(itr, maxItr)
		fmt.Println(*bm)
		fmt.Println("The average waiting time (opportunistic) is:", bm.AverageWaiting(maxItr))
		fmt.Println("The variance of the waiting time (opportunistic) is:", bm.VarianceWaiting(maxItr))
	*/

	itrSingleReq := 100
	itrReqs := 30
	maxItr := 10000
	averageNOPP := make([]float64, itrReqs)
	varianceNOPP := make([]float64, itrReqs)
	averageOPP := make([]float64, itrReqs)
	varianceOPP := make([]float64, itrReqs)
	bm := new(benchmark.Benchmarker)
	//bm.Set(itrSingleReq, "modified greedy", "grid")
	//bm.Set(itrSingleReq, "modified greedy", "ring")

	bm.Set(itrSingleReq, "nonoblivious local", "grid")
	//bm.Set(itrSingleReq, "nonoblivious local", "ring")

	//bm.Set(itrSingleReq, "qpass", "grid")
	//bm.Set(itrSingleReq, "qpass", "ring")
	bm.SetKeepReqs(true)
	for i := 0; i <= itrReqs-1; i++ {
		fmt.Println("Average Run:", i)
		bm.RegenerateReqs(itrSingleReq)
		config.SetOpportunism(false)
		bm.Start(itrSingleReq, maxItr)
		averageNOPP[i] = bm.AverageWaiting(maxItr)
		varianceNOPP[i] = bm.VarianceWaiting(maxItr)
		//fmt.Println(*bm)
		fmt.Println("NOPP Finished.")
		config.SetOpportunism(true)
		bm.Start(itrSingleReq, maxItr)
		fmt.Println(*bm)
		averageOPP[i] = bm.AverageWaiting(maxItr)
		varianceOPP[i] = bm.VarianceWaiting(maxItr)
	}
	fmt.Println("Average NOPP waiting time is:", AverageWaiting(averageNOPP, maxItr))
	fmt.Println("Average OPP waiting time is:", AverageWaiting(averageOPP, maxItr))
}

func AverageWaiting(nums []float64, maxItr int) float64 {
	sum := float64(0)
	meanLength := len(nums)
	for _, val := range nums {
		if val >= float64(maxItr-1) {
			meanLength -= 1
			continue
		}
		sum += val
	}
	return float64(sum) / float64(meanLength)
}
