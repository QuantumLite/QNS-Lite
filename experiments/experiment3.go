package main

import (
	"fmt"

	"example.com/benchmark"
	"example.com/config"
)

func experiment3(bm *benchmark.Benchmarker, itrReqs int, itrSingleReq int, maxItr int) {
	averageNOPP := make([]float64, itrReqs)
	varianceNOPP := make([]float64, itrReqs)
	averageOPP := make([]float64, itrReqs)
	varianceOPP := make([]float64, itrReqs)
	sizes := [11]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	ATWT := make([][]float64, 6)
	for i := 0; i <= 5; i++ {
		ATWT[i] = make([]float64, len(sizes))
	}
	algos := [3]string{"modified greedy", "nonoblivious local", "qpass"}
	topologies := [1]string{"grid"}
	config.SetPGen(0.8)
	config.SetPSwap(float64(0.8))
	config.SetLifetime(30)
	config.SetNumRequests(20)
	for sizeIndex, size := range sizes {
		config.SetSize(size)
		fmt.Println("size is", size)
		fmt.Println("config.size is", config.GetConfig().GetSize())
		for algo := 0; algo < 3; algo++ {
			fmt.Println("algorithm is", algos[algo])
			bm.Set(itrSingleReq, algos[algo], topologies[0])
			bm.SetKeepReqs(true)
			for i := 0; i <= itrReqs-1; i++ {
				//fmt.Println("Average Run:", i)
				bm.RegenerateReqs(itrSingleReq)
				config.SetOpportunism(false)
				bm.Start(itrSingleReq, maxItr)
				averageNOPP[i] = bm.AverageWaiting(maxItr)
				varianceNOPP[i] = bm.VarianceWaiting(maxItr)
				//fmt.Println(*bm)
				//fmt.Println("NOPP Finished.")
				config.SetOpportunism(true)
				bm.Start(itrSingleReq, maxItr)
				//fmt.Println(*bm)
				averageOPP[i] = bm.AverageWaiting(maxItr)
				varianceOPP[i] = bm.VarianceWaiting(maxItr)
			}
			ATWT[2*algo][sizeIndex] = AverageWaiting(averageNOPP, maxItr)
			ATWT[2*algo+1][sizeIndex] = AverageWaiting(averageOPP, maxItr)
			fmt.Println("Average NOPP waiting time is:", AverageWaiting(averageNOPP, maxItr))
			fmt.Println("Average OPP waiting time is:", AverageWaiting(averageOPP, maxItr))
		}
	}
	//file, err := os.OpenFile("./Data/experiment1.txt", os.O_APPEND|os.O_WRONLY, 0644)
	handleFile(ATWT, "./Data/experiment3.txt")
	/*var err = os.Remove("./Data/experiment1.txt")
	if err != nil {
		log.Println(err)
	}
	var _, errStat = os.Stat("./Data/experiment1.txt")

	// create file if not exists
	if os.IsNotExist(errStat) {
		var file, errCreate = os.Create("./Data/experiment1.txt")
		if errCreate != nil {
			log.Println(errCreate)
		}
		defer file.Close()
	}
	file, err := os.OpenFile("./Data/experiment1.txt", os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	ATWT_String := make([][]string, 6)
	for i := 0; i <= 5; i++ {
		ATWT_String[i] = make([]string, len(p))
	}
	for i := 0; i < 6; i++ {
		if _, err := file.WriteString("\n["); err != nil {
			log.Fatal(err)
		}
		for j := 0; j < len(p); j++ {
			ATWT_String[i][j] = strconv.Itoa(int(ATWT[i][j]))
			if j == len(p)-1 {
				if _, err := file.WriteString(ATWT_String[i][j]); err != nil {
					log.Fatal(err)
				}
			} else {
				if _, err := file.WriteString(ATWT_String[i][j] + ", "); err != nil {
					log.Fatal(err)
				}
			}
		}
		if _, err := file.WriteString("]"); err != nil {
			log.Fatal(err)
		}
	}*/

	fmt.Println(ATWT)
}
