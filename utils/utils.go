package utils

import (
	"example.com/graph"
)

type anything interface{}

func RemoveNode(slice []*graph.Node, index int) []*graph.Node {
	if index == 0 {
		return slice[1:len(slice)]
	} else if index == len(slice)-1 {
		return slice[0 : len(slice)-1]
	} else {
		temp := slice[0:index]
		tempEnd := slice[index+1 : len(slice)]
		temp = append(temp, tempEnd...)
		return temp
	}
}
