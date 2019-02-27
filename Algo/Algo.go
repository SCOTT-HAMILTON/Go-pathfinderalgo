package Algo

import "github.com/SCOTT-HAMILTON/Go-pathfinderalgo/Node"

type AlgoInterface interface {
	Init()
	Update(done chan bool)
	UpdateFinalPath()
}

type Algo struct {
	Nbw       int
	Nbh       int
	Start     int
	End       int
	CurNode   Node.NodeCase
	Map_walls *[]int
	Path      []int
	Finished  bool
}

func NewAlgo(nbw, nbh, s, e int, mapwalls *[]int) Algo {
	return Algo{nbw, nbh, s, e, Node.NewNodeCase(), mapwalls, make([]int, 0), false}
}
