package PathFinderAlgo

import "math"
import "github.com/SCOTT-HAMILTON/Go-pathfinderalgo/Node"

type Djikstra struct {
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

func NewDjikstra(nbw int, nbh int, s int, e int, mapwalls *[]int) Djikstra {
	return Djikstra{nbw, nbh, s, e, Node.NewNodeCase(), mapwalls, make([]int, 0), make([]Node.NodeCase, 0), make([]Node.NodeCase, 0), false}
}

func (djk *Djikstra) IsFinished() bool {
	return djk.finished
}

func (djk *Djikstra) Init() {
	djk.finished = false
	//Reset map cache
	for i, n := range *djk.Map_walls {
		if n != 1 {
			(*djk.Map_walls)[i] = 0
		}
	}

	//reset slices, keep memory allocated
	djk.openList = djk.openList[:0]
	djk.closedList = djk.closedList[:0]

	//Put to start
	djk.CurNode = Node.NodeCase{djk.Start, 0, 0, 0, 0}
	djk.openList = append(djk.openList, djk.CurNode)

	print("Djikstra initialized with sizes : ", djk.Nbw, "X", djk.Nbh, "\n")
}

func (djk *Djikstra) FindNei(pos int) Node.NodeCase {
	if pos == djk.CurNode.Pos {
		return djk.CurNode
	}
	for _, n := range djk.openList {
		if n.Pos == pos {
			return n
		}
	}
	for _, n := range djk.closedList {
		if n.Pos == pos {
			return n
		}
	}
	return Node.NodeCase{-1, 0, 0, 0, 0}
}

func (djk *Djikstra) IsInMap(pos int) bool {
	if pos >= 0 && pos < len(*djk.Map_walls) {
		return true
	} else {
		return false
	}
}

func (djk *Djikstra) IsWall(pos int) bool {
	return (*djk.Map_walls)[pos] == 1
}

func (djk *Djikstra) IsLeftBorder(pos int) bool {
	return pos%djk.Nbw == 0
}

func (djk *Djikstra) IsRightBorder(pos int) bool {
	return (pos+1)%djk.Nbw == 0
}

func (djk *Djikstra) IsTopBorder(pos int) bool {
	return pos < djk.Nbw
}

func (djk *Djikstra) IsBottomBorder(pos int) bool {
	return pos > (djk.Nbw*djk.Nbh)-1-djk.Nbw
}

func (djk *Djikstra) getNeighbours() (neighbours []Node.NodeCase) {
	tmpNei := make([]Node.Case, 8)
	nbNei := 0

	// left case
	leftMiddle := djk.CurNode.Pos - 1
	if djk.IsInMap(leftMiddle) && !djk.IsWall(leftMiddle) &&
		!djk.IsLeftBorder(djk.CurNode.Pos) {
		tmpNei[0] = Node.Case{Pos: leftMiddle, Val: true}
		nbNei++
	} else {
		tmpNei[0] = Node.Case{Pos: leftMiddle, Val: false}
	}

	// middle top case
	middleTop := djk.CurNode.Pos - djk.Nbw
	if djk.IsInMap(middleTop) && !djk.IsWall(middleTop) &&
		!djk.IsTopBorder(djk.CurNode.Pos) {
		tmpNei[1] = Node.Case{Pos: middleTop, Val: true}
		nbNei++
	} else {
		tmpNei[1] = Node.Case{Pos: middleTop, Val: false}
	}

	// right case
	rightMiddle := djk.CurNode.Pos + 1
	if djk.IsInMap(rightMiddle) && !djk.IsWall(rightMiddle) &&
		!djk.IsRightBorder(djk.CurNode.Pos) {
		tmpNei[2] = Node.Case{Pos: rightMiddle, Val: true}
		nbNei++
	} else {
		tmpNei[2] = Node.Case{Pos: rightMiddle, Val: false}
	}

	// middle bottom case
	middleBottom := djk.CurNode.Pos + djk.Nbw
	if djk.IsInMap(middleBottom) && !djk.IsWall(middleBottom) &&
		!djk.IsBottomBorder(djk.CurNode.Pos) {
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

func (djk *Djikstra) bestNode() (index int) {
	index = -1
	fB := float64(-1)
	for i, n := range djk.openList {
		if fB == -1 || n.F < fB || (fB == n.F && n.G < djk.openList[index].G) {
			if Node.FindNei(&djk.closedList, &n) == -1 {
				index = i
				fB = n.F
			}
		}
	}
	return
}

func (djk *Djikstra) ToCoord(pos int) (x int, y int) {
	x = pos % djk.Nbw
	y = (pos - x) / djk.Nbw
	return
}

func (djk *Djikstra) Dist(nodeposA int, nodeposB int) (dist float64) {
	type pos struct {
		x float64
		y float64
	}
	A := pos{0, 0}
	{
		x, y := djk.ToCoord(nodeposA)
		A.x = float64(x)
		A.y = float64(y)
	}

	B := pos{0, 0}
	{
		x, y := djk.ToCoord(nodeposB)
		B.x = float64(x)
		B.y = float64(y)
	}
	t := (B.x - A.x)
	v := (B.y - A.y)
	dist = math.Sqrt(t*t + v*v)
	return
}

func (djk *Djikstra) Update(done chan bool) {
	defer func() {
		done <- true
	}()

	if djk.finished {
		return
	}

	//get Best Node
	index := djk.bestNode()
	djk.CurNode = djk.openList[index]
	//print("best : ", djk.CurNode.Pos,"\n")

	//remove it from openList
	djk.openList[index] = djk.openList[len(djk.openList)-1]
	djk.openList = djk.openList[:len(djk.openList)-1]

	//add it to closedList
	djk.closedList = append(djk.closedList, djk.CurNode)

	//is the end
	if djk.CurNode.Pos == djk.End {
		djk.finished = true
		return
	}

	neighbours := djk.getNeighbours()

	for i, n := range *djk.Map_walls {
		if n != 1 {
			(*djk.Map_walls)[i] = 0
		}
	}

	for _, nei := range neighbours {
		if Node.FindNei(&djk.closedList, &nei) != -1 {
			continue
		}

		(*djk.Map_walls)[nei.Pos] = 2

		nei.G = djk.CurNode.G + 1
		nei.F = nei.G
		nei.Parent = djk.CurNode.Pos
		stop := false
		for _, n := range djk.openList {
			if n.Pos == nei.Pos && nei.G > n.G {
				//stop her
				stop = true
				continue
			}
		}
		if !stop {
			djk.openList = append(djk.openList, nei)
		}

	}

	//update openlist to map
	for _, n := range djk.openList {
		if (*djk.Map_walls)[n.Pos] != 2 { //is not already neighbour
			(*djk.Map_walls)[n.Pos] = 3
		}
	}

	//update closedList to map
	for _, n := range djk.closedList {
		if (*djk.Map_walls)[n.Pos] != 2 { //is not already neighbour
			(*djk.Map_walls)[n.Pos] = 4
		}
	}
}

func (djk *Djikstra) UpdateFinalPath() {
	cur := djk.CurNode.Pos
	djk.Path = djk.Path[:0]
	for cur != djk.Start {
		djk.Path = append(djk.Path, cur)
		node := djk.FindNei(cur)
		if node.Pos == -1 {
			print("error no node at pos : ", cur)
			return
		}
		cur = node.Parent
	}
}
