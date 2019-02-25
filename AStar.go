package PathFinderAlgo

import "math"
import "github.com/SCOTT-HAMILTON/PathFinderAlgo/Node"

type AStar struct {
	Nbw        int
	Nbh        int
	Start      int
	End        int
	CurNode    Node.NodeCase
	Map_walls  *[]int
	Path       []int
	openList   []Node.NodeCase
	closedList []Node.NodeCase
	finished   bool
}

func NewAStar(nbw int, nbh int, s int, e int, mapwalls *[]int) AStar {
	return AStar{nbw, nbh, s, e, Node.NewNodeCase(), mapwalls, make([]int, 0), make([]Node.NodeCase, 0), make([]Node.NodeCase, 0), false}
}

func (star *AStar) IsFinished() bool {
	return star.finished
}

func (star *AStar) Init() {
	star.finished = false
	//Reset map cache
	for i, n := range *star.Map_walls {
		if n != 1 {
			(*star.Map_walls)[i] = 0
		}
	}

	//reset slices, keep memory allocated
	star.openList = star.openList[:0]
	star.closedList = star.closedList[:0]

	//Put to start
	star.CurNode = Node.NodeCase{star.Start, 0, 0, 0, 0}
	star.openList = append(star.openList, star.CurNode)

	print("AStar initialized with sizes : ", star.Nbw, "X", star.Nbh, "\n")
}

func (star *AStar) FindNei(pos int) Node.NodeCase {
	if pos == star.CurNode.Pos {
		return star.CurNode
	}
	for _, n := range star.openList {
		if n.Pos == pos {
			return n
		}
	}
	for _, n := range star.closedList {
		if n.Pos == pos {
			return n
		}
	}
	return Node.NodeCase{-1, 0, 0, 0, 0}
}

func (star *AStar) IsInMap(pos int) bool {
	if pos >= 0 && pos < len(*star.Map_walls) {
		return true
	} else {
		return false
	}
}

func (star *AStar) IsWall(pos int) bool {
	return (*star.Map_walls)[pos] == 1
}

func (star *AStar) IsLeftBorder(pos int) bool {
	return pos%star.Nbw == 0
}

func (star *AStar) IsRightBorder(pos int) bool {
	return (pos+1)%star.Nbw == 0
}

func (star *AStar) IsTopBorder(pos int) bool {
	return pos < star.Nbw
}

func (star *AStar) IsBottomBorder(pos int) bool {
	return pos > (star.Nbw*star.Nbh)-1-star.Nbw
}

func (star *AStar) getNeighbours() (neighbours []Node.NodeCase) {
	tmpNei := make([]Node.Case, 8)
	nbNei := 0

	// left case
	leftMiddle := star.CurNode.Pos - 1
	if star.IsInMap(leftMiddle) && !star.IsWall(leftMiddle) &&
		!star.IsLeftBorder(star.CurNode.Pos) {
		tmpNei[0] = Node.Case{Pos: leftMiddle, Val: true}
		nbNei++
	} else {
		tmpNei[0] = Node.Case{Pos: leftMiddle, Val: false}
	}

	// middle top case
	middleTop := star.CurNode.Pos - star.Nbw
	if star.IsInMap(middleTop) && !star.IsWall(middleTop) &&
		!star.IsTopBorder(star.CurNode.Pos) {
		tmpNei[1] = Node.Case{Pos: middleTop, Val: true}
		nbNei++
	} else {
		tmpNei[1] = Node.Case{Pos: middleTop, Val: false}
	}

	// right case
	rightMiddle := star.CurNode.Pos + 1
	if star.IsInMap(rightMiddle) && !star.IsWall(rightMiddle) &&
		!star.IsRightBorder(star.CurNode.Pos) {
		tmpNei[2] = Node.Case{Pos: rightMiddle, Val: true}
		nbNei++
	} else {
		tmpNei[2] = Node.Case{Pos: rightMiddle, Val: false}
	}

	// middle bottom case
	middleBottom := star.CurNode.Pos + star.Nbw
	if star.IsInMap(middleBottom) && !star.IsWall(middleBottom) &&
		!star.IsBottomBorder(star.CurNode.Pos) {
		tmpNei[3] = Node.Case{Pos: middleBottom, Val: true}
		nbNei++
	} else {
		tmpNei[3] = Node.Case{Pos: middleBottom, Val: false}
	}

	neighbours = make([]Node.NodeCase, nbNei)
	i2 := 0
	for i, _ := range tmpNei {
		if tmpNei[i].Val {
			//			print("yes new neighbour!!")
			neighbours[i2] = Node.NodeCase{tmpNei[i].Pos, 0, 0, 0, 0}
			i2++
		}
	}

	return
}

func (star *AStar) bestNode() (index int) {
	index = -1
	fB := float64(-1)
	for i, n := range star.openList {
		if fB == -1 || n.F < fB || (fB == n.F && n.G < star.openList[index].G) {
			if Node.FindNei(&star.closedList, &n) == -1 {
				index = i
				fB = n.F
			}
		}
	}
	return
}

func (star *AStar) ToCoord(pos int) (x int, y int) {
	x = pos % star.Nbw
	y = (pos - x) / star.Nbw
	return
}

func (star *AStar) Dist(nodeposA int, nodeposB int) (dist float64) {
	type pos struct {
		x float64
		y float64
	}
	A := pos{0, 0}
	{
		x, y := star.ToCoord(nodeposA)
		A.x = float64(x)
		A.y = float64(y)
	}

	B := pos{0, 0}
	{
		x, y := star.ToCoord(nodeposB)
		B.x = float64(x)
		B.y = float64(y)
	}
	t := (B.x - A.x)
	v := (B.y - A.y)
	dist = math.Sqrt(t*t + v*v)
	return
}

func (star *AStar) Update(done chan bool) {
	defer func() {
		done <- true
	}()

	if star.finished {
		return
	}

	//get Best Node
	index := star.bestNode()
	star.CurNode = star.openList[index]
	//print("best : ", star.CurNode.Pos,"\n")

	//remove it from openList
	star.openList[index] = star.openList[len(star.openList)-1]
	star.openList = star.openList[:len(star.openList)-1]

	//add it to closedList
	star.closedList = append(star.closedList, star.CurNode)

	//is the end
	if star.CurNode.Pos == star.End {
		star.finished = true
		return
	}

	neighbours := star.getNeighbours()

	for i, n := range *star.Map_walls {
		if n != 1 {
			(*star.Map_walls)[i] = 0
		}
	}

	for _, nei := range neighbours {
		if Node.FindNei(&star.closedList, &nei) != -1 {
			continue
		}

		(*star.Map_walls)[nei.Pos] = 2

		nei.G = star.CurNode.G + 1
		nei.H = star.Dist(nei.Pos, star.End)
		nei.H *= nei.H
		nei.F = nei.G + nei.H
		nei.Parent = star.CurNode.Pos
		stop := false
		for _, n := range star.openList {
			if n.Pos == nei.Pos && nei.G > n.G {
				//stop her
				stop = true
				continue
			}
		}
		if !stop {
			star.openList = append(star.openList, nei)
		}

	}

	//update openlist to map
	for _, n := range star.openList {
		if (*star.Map_walls)[n.Pos] != 2 { //is not already neighbour
			(*star.Map_walls)[n.Pos] = 3
		}
	}

	//update closedList to map
	for _, n := range star.closedList {
		if (*star.Map_walls)[n.Pos] != 2 { //is not already neighbour
			(*star.Map_walls)[n.Pos] = 4
		}
	}
}

func (star *AStar) UpdateFinalPath() {
	cur := star.CurNode.Pos
	star.Path = star.Path[:0]
	for cur != star.Start {
		star.Path = append(star.Path, cur)
		node := star.FindNei(cur)
		if node.Pos == -1 {
			print("error no node at pos : ", cur)
			return
		}
		cur = node.Parent
	}
}
