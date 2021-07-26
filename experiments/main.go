package main

import (
	"log"
	"math"
	"os"
	"strconv"

	"example.com/benchmark"
)

func main() {
	itrSingleReq := 100
	itrReqs := 30
	maxItr := 10000
	bm := new(benchmark.Benchmarker)
	//experiment1(bm, itrReqs, itrSingleReq, maxItr)
	//experiment2(bm, itrReqs, itrSingleReq, maxItr)
	//experiment3(bm, itrReqs, itrSingleReq, maxItr)
	//experiment4(bm, itrReqs, itrSingleReq, maxItr)
	//experiment5(bm, itrReqs, itrSingleReq, maxItr)
	experiment6(bm, itrReqs, itrSingleReq, maxItr)
}

func AverageWaiting(nums []float64, maxItr int) float64 {
	sum := float64(0)
	meanLength := len(nums)
	for _, val := range nums {
		if val >= float64(maxItr-1) {
			meanLength -= 1
			continue
		}
		if math.IsNaN(val) {
			meanLength -= 1
			continue
		}
		sum += val
	}
	return float64(sum) / float64(meanLength)
}

func VarianceWaiting(nums []float64, maxItr int) float64 {
	sum := float64(0)
	ave := AverageWaiting(nums, maxItr)
	varLength := len(nums)
	for _, val := range nums {
		if val >= float64(maxItr)-1 {
			varLength--
			continue
		}
		if math.IsNaN(val) {
			varLength -= 1
			continue
		}
		sum += (float64(val) - ave) * (float64(val) - ave)
	}
	return float64(sum) / float64(varLength)
}

/*func (bm *Benchmarker) VarianceWaiting(maxItr int) float64 {
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
}*/

func handleFile(data [][]float64, filePath string) {
	var err = os.Remove(filePath)
	if err != nil {
		log.Println(err)
	}
	var _, errStat = os.Stat(filePath)

	// create file if not exists
	if os.IsNotExist(errStat) {
		var file, errCreate = os.Create(filePath)
		if errCreate != nil {
			log.Println(errCreate)
		}
		defer file.Close()
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	data_String := make([][]string, len(data))
	for i := 0; i <= len(data)-1; i++ {
		data_String[i] = make([]string, len(data[0]))
	}
	for i := 0; i < len(data); i++ {
		if _, err := file.WriteString("\n["); err != nil {
			log.Fatal(err)
		}
		for j := 0; j < len(data[0]); j++ {
			data_String[i][j] = strconv.Itoa(int(data[i][j]))
			if j == len(data[0])-1 {
				if _, err := file.WriteString(data_String[i][j]); err != nil {
					log.Fatal(err)
				}
			} else {
				if _, err := file.WriteString(data_String[i][j] + ", "); err != nil {
					log.Fatal(err)
				}
			}
		}
		if _, err := file.WriteString("]"); err != nil {
			log.Fatal(err)
		}
	}
}
