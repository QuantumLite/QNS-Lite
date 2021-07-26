package main

import (
	"fmt"

	"example.com/benchmark"
	"example.com/config"
)

func experiment5(bm *benchmark.Benchmarker, itrReqs int, itrSingleReq int, maxItr int) {
	averageNOPP := make([]float64, itrReqs)
	varianceNOPP := make([]float64, itrReqs)
	averageOPP := make([]float64, itrReqs)
	varianceOPP := make([]float64, itrReqs)
	priority1NOPP := make([]float64, itrReqs)
	priority1OPP := make([]float64, itrReqs)
	priority2NOPP := make([]float64, itrReqs)
	priority2OPP := make([]float64, itrReqs)
	priority3NOPP := make([]float64, itrReqs)
	priority3OPP := make([]float64, itrReqs)
	////// Make such things parametric!
	temp := make([]float64, 3)
	p := [10]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	ATWT := make([][]float64, 2)
	for i := 0; i <= 1; i++ {
		ATWT[i] = make([]float64, len(p))
	}
	PATWT := make([][]float64, 6)
	for i := 0; i <= 5; i++ {
		PATWT[i] = make([]float64, len(p))
	}
	algos := [1]string{"modified greedy"}
	topologies := [1]string{"grid"}
	config.SetPSwap(float64(1))
	config.SetSize(10)
	config.SetLifetime(22)
	config.SetNumRequests(30)
	for p_genIndex, p_gen := range p {
		config.SetPGen(p_gen)
		fmt.Println("p_gen is", p_gen)
		fmt.Println("config.p_gen is", config.GetConfig().GetPGen())
		for algo := 0; algo < len(algos); algo++ {
			fmt.Println("algorithm is", algos[algo])
			bm.Set(itrSingleReq, algos[algo], topologies[0])
			bm.SetKeepReqs(true)
			bm.RegenerateReqs(itrSingleReq)
			for i := 0; i <= itrReqs-1; i++ {
				//fmt.Println("Average Run:", i)
				//bm.RegenerateReqs(itrSingleReq)
				config.SetOpportunism(false)
				bm.Start(itrSingleReq, maxItr)
				averageNOPP[i] = bm.AverageWaiting(maxItr)
				varianceNOPP[i] = bm.VarianceWaiting(maxItr)
				temp = bm.PriorityAverageWaiting(maxItr)
				priority1NOPP[i] = temp[0]
				priority2NOPP[i] = temp[1]
				priority3NOPP[i] = temp[2]
				//fmt.Println(*bm)
				//fmt.Println("NOPP Finished.")
				config.SetOpportunism(true)
				bm.Start(itrSingleReq, maxItr)
				//fmt.Println(*bm)
				averageOPP[i] = bm.AverageWaiting(maxItr)
				varianceOPP[i] = bm.VarianceWaiting(maxItr)
				temp = bm.PriorityAverageWaiting(maxItr)
				priority1OPP[i] = temp[0]
				priority2OPP[i] = temp[1]
				priority3OPP[i] = temp[2]
			}
			ATWT[2*algo][p_genIndex] = AverageWaiting(averageNOPP, maxItr)
			ATWT[2*algo+1][p_genIndex] = AverageWaiting(averageOPP, maxItr)
			fmt.Println("Average NOPP waiting time is:", AverageWaiting(averageNOPP, maxItr))
			fmt.Println("Average OPP waiting time is:", AverageWaiting(averageOPP, maxItr))
			PATWT[0][p_genIndex] = AverageWaiting(priority1NOPP, maxItr)
			PATWT[1][p_genIndex] = AverageWaiting(priority2NOPP, maxItr)
			PATWT[2][p_genIndex] = AverageWaiting(priority3NOPP, maxItr)
			PATWT[3][p_genIndex] = AverageWaiting(priority1OPP, maxItr)
			PATWT[4][p_genIndex] = AverageWaiting(priority2OPP, maxItr)
			PATWT[5][p_genIndex] = AverageWaiting(priority3OPP, maxItr)
		}
	}
	//file, err := os.OpenFile("./Data/experiment1.txt", os.O_APPEND|os.O_WRONLY, 0644)
	handleFile(PATWT, "./Data/experiment5.txt")
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
