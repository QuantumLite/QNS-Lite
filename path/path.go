package path

import (
	"fmt"
	"math"

	"example.com/config"
	"example.com/graph"
	"example.com/request"
)

const (
	MODIFIED_GREEDY    = "modified greedy"
	NONOBLIVIOUS_LOCAL = "nonoblivious local"
	Q_PASS             = "Q-Pass"
	Q_CAST             = "Q-CAST"
	SLMP               = "SLMP"
)

type Path []*graph.Node

// The PathFinder interface captures the path finding algorithm abstraction.
type PathFinder interface {
	Build(graph.Topology)
	Clear()
	Find(src, dest *graph.Node) (Path, []int, []int, bool)
	//GetPath() path
	next(dest *graph.Node) (*graph.Node, []*graph.Node, int)
	add(*graph.Node)
	//getNetwork() *graph.Topology
}

func BuildPathFinder(algorithm string, network graph.Topology) PathFinder {
	if algorithm == MODIFIED_GREEDY {
		var mg *modifiedGreedy = new(modifiedGreedy)
		mg.Build(network)
		return mg
	} else if algorithm == NONOBLIVIOUS_LOCAL {
		var nol *nonObliviousLocal = new(nonObliviousLocal)
		nol.Build(network)
		return nol
	}
	fmt.Println("path.go: Warning! algorithm not recognized!")
	return nil
}

func PF(network graph.Topology, reqs []*request.Request, algorithm string, ignoreLeftOvers bool) {
	pf := BuildPathFinder(algorithm, network)
	hasContention := config.GetConfig().GetHasContention()
	leftOverReqs := make([]*request.Request, 0)
	//paths := make([]Path, len(reqs))
	for i, req := range reqs {
		aggressiveness := config.GetConfig().GetAggressiveness()
		if !hasContention {
			//fmt.Println("Hi, I do not have contention!")
			for j := 1; j <= aggressiveness; j++ {
				//fmt.Println("j is", j)
				if j == 1 {
					//fmt.Println("Hi, I do not have contention!, j = 1")
					nodes, _, _, _ := pf.Find(req.Src, req.Dest)
					//if len(nodes) == 0 {
					//	continue
					//}
					nodesToBeCopied := PathToNode(nodes)
					reqs[i].Paths[j-1] = make([]*graph.Node, len(nodesToBeCopied))
					copy(reqs[i].Paths[j-1], nodesToBeCopied)
				} else {
					//fmt.Println("Hi, I do not have contention!, j = ", j)
					nodes, _, _, _ := pf.Find(req.Src, req.Dest)
					//if len(nodes) == 0 {
					//	continue
					//}
					reqs[i].Paths = append(reqs[i].Paths, PathToNode(nodes))
				}
				//fmt.Println("Hi, I do not have contention! Out of j if-else.")
				graph.Prune(PathToLinks(reqs[i].Paths[len(reqs[i].Paths)-1], network))
				//fmt.Println("Hi, I do not have contention! Just pruned!")
				pf.Clear()
			}
			//fmt.Println("Out of j for")
			if len(req.Paths[0]) == 0 {
				leftOverReqs = append(leftOverReqs, req)
			}
			if len(req.Paths) > 1 {
				for m, _ := range req.Paths {
					if m == 0 {
						continue
					}
					//fmt.Println("req.Path length is:", len(req.Paths))
					if len(req.Paths[m]) == 0 {
						req.Paths = req.Paths[0:m]
						break
					}
				}
			}
			//fmt.Println("before continue")
			continue
		}
		//fmt.Println("out of contention")
		//fmt.Println("PF - request number", i)
		nodes, mapping, options, mappingIsNull := pf.Find(req.Src, req.Dest)
		pf.Clear()
		if len(nodes) == 0 {
			//fmt.Println("Null path found.")
		}
		//fmt.Println("mapping is", mapping)
		//fmt.Println("first path is")

		//for _, nn := range nodes {
		//	fmt.Println("node", nn.ID)
		//}
		nodesToBeCopied := PathToNode(nodes)
		reqs[i].Paths[0] = make([]*graph.Node, len(nodesToBeCopied))
		copy(reqs[i].Paths[0], nodesToBeCopied)
		if aggressiveness == 1 {
			if !hasContention {
				graph.Prune(PathToLinks(reqs[i].Paths[len(reqs[i].Paths)-1], network))
			}
			continue
		} else {
			if mappingIsNull {
				continue
			}
			totalOptions := sum(options)
			totalOptions = int(math.Min(float64(totalOptions), float64(aggressiveness)-1))
			aggIndex := 0
			total := 0
			/*tempIndex := 1
			index := 0
			tries := 1
			linksToDeprune := make([]*graph.Link, 0)
			var pathToBeCopied []*graph.Node*/

			//tempLink := make([]*graph.Link, 1)
			//tempLink[0] = network.GetLinkBetween(reqs[i].Paths[aggIndex][mapping[0]], reqs[i].Paths[total][mapping[0]+1])
			//graph.Prune(tempLink)

			////////////////////////////////////// Useful logging
			//fmt.Println("mapping[index] is", mapping[0])
			//fmt.Println("Source of the request is", req.Src.ID, "Destination of the request is", req.Dest.ID)
			//fmt.Println("totalOptions is", totalOptions)

			findOverlappingPaths(reqs, network, pf, &total, &totalOptions, aggIndex, 0, i, mapping, options)
			/*for tries <= totalOptions && total <= totalOptions {
				if index >= len(mapping) {
					tries++
					continue
				}
				fmt.Println("tries is", tries, "total is", total)
				// Should I use Paths[index] or Paths[0]?????????????????????????????????????????????????????
				// I should use all of the paths for further exploration, not just Paths[0]
				tries++
				tempLink := make([]*graph.Link, 1)
				fmt.Println("mapping[index] is", mapping[index])
				tempLink[0] = network.GetLinkBetween(reqs[i].Paths[aggIndex][mapping[index]], reqs[i].Paths[total][mapping[index]+1])
				graph.Prune(tempLink)
				linksToDeprune = append(linksToDeprune, tempLink...)
				tail, _, _, _ := pf.Find(reqs[i].Paths[aggIndex][mapping[index]], req.Dest)
				pf.Clear()
				if tail != nil {
					fmt.Println("Aha! Found one for request number", i)
					temp := make([]*graph.Node, len(reqs[i].Paths[aggIndex][0:mapping[index]+1])+len(tail)-1)
					reqs[i].Paths = append(reqs[i].Paths, temp)
					//reqs[i].Paths[total+1] = make([]*graph.Node, len(reqs[i].Paths[aggIndex][0:mapping[index]+1])+len(tail))
					if mapping[index] == 0 {
						pathToBeCopied = PathToNode(tail)
					} else {
						///////////////////////////////////////// Check the 0:mapping[index]
						pathToBeCopied = append(reqs[i].Paths[aggIndex][0:mapping[index]], PathToNode(tail)...)
					}
					copy(reqs[i].Paths[total+1], pathToBeCopied)
					total++
				}
				tempIndex++
				if tempIndex == options[index] {
					fmt.Println("Ready to move to the next mapping node, and deprune the links.")
					for _, linkk := range linksToDeprune {
						fmt.Println("Depruning link", linkk.ID)
					}
					index++
					graph.DepruneLinks(linksToDeprune)
					linksToDeprune = make([]*graph.Link, 0)
					tempIndex = 0
				}
			}*/
		}

		//for j := 1; j <= config.GetConfig().GetAggressiveness(); j++ {
		//	if j == 1 {
		//		nodes := PathToNode(pf.Find(req.Src, req.Dest))
		//		reqs[i].Paths[j-1] = make([]*graph.Node, len(nodes))
		//		copy(reqs[i].Paths[j-1], nodes)
		//	} else {
		//		reqs[i].Paths = append(reqs[i].Paths, PathToNode(pf.Find(req.Src, req.Dest)))
		//	}
		//	if !hasContention {
		//		////////////////////////////// COMPLETE THIS!!!!!!!!!!!!!!!!!!!!!!
		//		//fmt.Println("En route to PathToLinks", reqs[i].Paths[len(reqs[i].Paths)-1])
		//		graph.Prune(PathToLinks(reqs[i].Paths[len(reqs[i].Paths)-1], network))
		//	}
		//	pf.Clear()
		//}
	}
	graph.Deprune(network)
	if ignoreLeftOvers {
		return
	}
	leftOverCntr := 0
	//fmt.Println("The length of the leftover requests is:", len(leftOverReqs))
	for len(leftOverReqs) > 0 {
		//fmt.Println("Dealing with leftover requests.")
		leftOverCntr++
		if leftOverCntr >= len(leftOverReqs)*len(leftOverReqs) {
			graph.Deprune(network)
			leftOverCntr = 0
		}
		temp := leftOverReqs
		leftOverReqs = make([]*request.Request, 0)
		for _, req := range temp {
			//fmt.Println("Hi, I do not have contention!, j = 1")
			nodes, _, _, _ := pf.Find(req.Src, req.Dest)
			pf.Clear()
			nodesToBeCopied := PathToNode(nodes)
			req.Paths[0] = make([]*graph.Node, len(nodesToBeCopied))
			copy(req.Paths[0], nodesToBeCopied)
			if len(req.Paths[0]) == 0 {
				leftOverReqs = append(leftOverReqs, req)
				continue
			}
			graph.Prune(PathToLinks(req.Paths[0], network))
		}
	}
	graph.Deprune(network)
	//return paths
}

func findOverlappingPaths(reqs []*request.Request, network graph.Topology, pf PathFinder, _total *int, totalOptions *int, aggIndex int, index int, i int, mapping []int, options []int) {
	//fmt.Println("Halo! Ich bin Ali!")
	linksToDeprune := make([]*graph.Link, 0)
	if len(mapping) > index+1 {
		//fmt.Println("Going one level deeper. index is", index)
		findOverlappingPaths(reqs, network, pf, _total, totalOptions, aggIndex, index+1, i, mapping, options)
		// We need some pruning here.
	}
	////////////// If found enough, just return.
	if *_total == *totalOptions {
		return
	}
	tempLink := make([]*graph.Link, 1)
	//fmt.Println("Right in the entrance:", "mapping[index] is", mapping[index], "aggIndex is", aggIndex, "total is", *_total)
	tempLink[0] = network.GetLinkBetween(reqs[i].Paths[aggIndex][mapping[index]], reqs[i].Paths[aggIndex][mapping[index]+1])
	graph.Prune(tempLink)
	linksToDeprune = append(linksToDeprune, tempLink...)

	//for _, linkkk := range linksToDeprune {
	//fmt.Println("link to deprune", linkkk.ID, "request number", i)
	//}

	//index := 0
	//aggIndex := 0
	///////////////////////////////////////////// tempIndex is important!!!!!
	tempIndex := 1
	tries := 1

	////////////////////////////////////////////////////// We need a local total

	//tries := *_tries
	req := reqs[i]
	//total := 0
	//total := *_total
	var pathToBeCopied []*graph.Node
	//fmt.Println("Source of the request is", req.Src.ID, "Destination of the request is", req.Dest.ID)
	//fmt.Println("totalOptions is", *totalOptions)
	////////////////////////////////////////////// Check the third condition
	for tries <= *totalOptions && *_total <= *totalOptions && tempIndex < options[index] {
		if index >= len(mapping) {
			tries++
			continue
		}
		//fmt.Println("tries is", tries, "total is", *_total)
		// Should I use Paths[index] or Paths[0]?????????????????????????????????????????????????????
		// I should use all of the paths for further exploration, not just Paths[0]
		tries++
		tempIndex++
		tail, newMapping, newOptions, newMappingIsNil := pf.Find(reqs[i].Paths[aggIndex][mapping[index]], req.Dest)
		*totalOptions += sum(newOptions)
		*totalOptions = int(math.Min(float64(*totalOptions), float64(config.GetConfig().GetAggressiveness()-1)))

		// Pruning ------------> Deprune when returning from the function!!!
		//tempLink := make([]*graph.Link, 1)
		//fmt.Println("mapping[index] is", mapping[index], "aggIndex is", aggIndex, "total is", *_total)
		//tempLink[0] = network.GetLinkBetween(reqs[i].Paths[aggIndex][mapping[index]], reqs[i].Paths[localTotal][mapping[index]+1])
		//graph.Prune(tempLink)
		//linksToDeprune = append(linksToDeprune, tempLink...)

		//////////// Previous find
		//fmt.Println("totalOptions is", *totalOptions)
		pf.Clear()
		if tail != nil {
			///////////////// Adding the new Path.
			//fmt.Println("Aha! New tail is not nil for request", i)
			temp := make([]*graph.Node, len(reqs[i].Paths[aggIndex][0:mapping[index]+1])+len(tail)-1)
			reqs[i].Paths = append(reqs[i].Paths, temp)
			//reqs[i].Paths[total+1] = make([]*graph.Node, len(reqs[i].Paths[aggIndex][0:mapping[index]+1])+len(tail))
			if mapping[index] == 0 {
				pathToBeCopied = PathToNode(tail)
			} else {
				///////////////////////////////////////// Check the 0:mapping[index]
				tempPath := make([]*graph.Node, 0)
				tempPath = append(tempPath, reqs[i].Paths[aggIndex][0:mapping[index]]...)
				pathToBeCopied = append(tempPath, PathToNode(tail)...)
				//pathToBeCopied = append(reqs[i].Paths[aggIndex][0:mapping[index]], PathToNode(tail)...)
			}
			//fmt.Println("Adding new path:", "*_total is", *_total)
			copy(reqs[i].Paths[*_total+1], pathToBeCopied)

			*_total++
			localTotal := *_total
			if !newMappingIsNil {
				//fmt.Println("Aha! New mapping is not nil for request", i, "The mapping is", newMapping)
				//////// Update the mapping indices.
				for j, _ := range newMapping {
					newMapping[j] += mapping[index]
				}
				//////// If enough paths are found, just return.
				if *_total == *totalOptions {
					graph.DepruneLinks(linksToDeprune)
					return
				}
				//graph.DepruneLinks(linksToDeprune)

				//for _, linkkk := range linksToDeprune {
				//	//fmt.Println("link to deprune", linkkk.ID, "request number", i)
				//}

				findOverlappingPaths(reqs, network, pf, _total, totalOptions, localTotal, 0, i, newMapping, newOptions)
				if *_total == *totalOptions {
					graph.DepruneLinks(linksToDeprune)
					//fmt.Println("Returning out!!!")
					return
				}
				//graph.Prune(linksToDeprune)
				///////////////////////////////////////// Why have I used -- here?!!
				//*_total--
				//tempIndex++

				//fmt.Println("tempIndex is", tempIndex, "options[index] is", options[index], "mapping[index] is", mapping[index])

				//if tempIndex == options[index] {
				//	fmt.Println("Ready to move to the next mapping node, and deprune the links.")
				//	for _, linkk := range linksToDeprune {
				//		fmt.Println("Depruning link", linkk.ID)
				//	}
				//	index++
				//	graph.DepruneLinks(linksToDeprune)
				//	linksToDeprune = make([]*graph.Link, 0)
				//	tempIndex = 0
				//}
				//continue
			}
			////////// Pruning
			tempLink := make([]*graph.Link, 1)
			//fmt.Println("mapping[index] is", mapping[index], "aggIndex is", aggIndex, "total is", *_total)
			tempLink[0] = network.GetLinkBetween(reqs[i].Paths[aggIndex][mapping[index]], reqs[i].Paths[localTotal][mapping[index]+1])
			graph.Prune(tempLink)
			linksToDeprune = append(linksToDeprune, tempLink...)

			//fmt.Println("Aha! Found one for request number", i)
			//temp := make([]*graph.Node, len(reqs[i].Paths[aggIndex][0:mapping[index]+1])+len(tail)-1)
			//reqs[i].Paths = append(reqs[i].Paths, temp)
			//reqs[i].Paths[total+1] = make([]*graph.Node, len(reqs[i].Paths[aggIndex][0:mapping[index]+1])+len(tail))
			//if mapping[index] == 0 {
			//	pathToBeCopied = PathToNode(tail)
			//} else {
			///////////////////////////////////////// Check the 0:mapping[index]
			//	pathToBeCopied = append(reqs[i].Paths[aggIndex][0:mapping[index]], PathToNode(tail)...)
			//}
			//copy(reqs[i].Paths[*_total+1], pathToBeCopied)
			//*_total++
		}
		//tempIndex++

		//if tempIndex == options[index] {
		//	fmt.Println("Ready to move to the next mapping node, and deprune the links.")
		//	for _, linkk := range linksToDeprune {
		//		fmt.Println("Depruning link", linkk.ID)
		//	}
		//	index++
		//	graph.DepruneLinks(linksToDeprune)
		//	linksToDeprune = make([]*graph.Link, 0)
		//	tempIndex = 0
		//}
	}
	//fmt.Println("Came out of the for loop.")
	graph.DepruneLinks(linksToDeprune)
}

//func batchFind()

//func copyNetwork(network graph.Topology) copied graph.Topology

// TODO: Check this!
func PathToNode(path Path) []*graph.Node {
	if path == nil {
		return nil
	}
	//fmt.Println("Inside PathToNodes. Input:", path)
	nodes := make([]*graph.Node, len(path))
	copy(nodes, path)
	//fmt.Println("Inside PathToNodes. Output:", nodes)
	return nodes
}

func PathToLinks(path Path, network graph.Topology) []*graph.Link {
	links := make([]*graph.Link, 0)
	i := 0
	//fmt.Println("Inside PathToLinks. Length of path:", len(path), "Path is:", path)
	//for _, node := range path {
	//	fmt.Println("NODE", node.ID)
	//}
	for i <= len(path)-2 {
		links = append(links, network.GetLinkBetween(path[i], path[i+1]))
		i++
	}
	return links
}

func sum(in []int) int {
	s := 0
	for _, val := range in {
		s += val
	}
	return s
}
