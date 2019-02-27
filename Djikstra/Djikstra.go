package Djikstra

import "github.com/SCOTT-HAMILTON/Go-pathfinderalgo/Node"
import "github.com/SCOTT-HAMILTON/Go-pathfinderalgo/AStar"

type Djikstra struct {
	*AStar.AStar
}

func NewDjikstra(nbw int, nbh int, s int, e int, mapwalls *[]int) *Djikstra {
	djk := &Djikstra{AStar: AStar.NewAStar(nbw, nbh, s, e, mapwalls)}
	djk.FctSetupNei = djk.setupNeiImpl
	return djk
}

func (djk *Djikstra) setupNeiImpl(nei *Node.NodeCase) {
	nei.G = djk.CurNode.G + 1
	nei.F = nei.G
	nei.Parent = djk.CurNode.Pos
}
