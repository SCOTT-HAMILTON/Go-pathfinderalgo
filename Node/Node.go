package Node

type NodeCase struct {
	Pos    int
	Parent int
	G      float64
	H      float64
	F      float64
}

func NewNodeCase() NodeCase {
	return NodeCase{0, 0, 0, 0, 0}
}

func (node *NodeCase) GetF() float64 {
	return node.F
}

func (node *NodeCase) GetG() float64 {
	return node.G
}

func (node *NodeCase) GetH() float64 {
	return node.H
}

func FindNei(tab *[]NodeCase, nei *NodeCase) int {
	for i, n := range *tab {
		if n.Pos == nei.Pos {
			return i
		}
	}
	return -1
}

type Case struct {
	Pos int
	Val bool
}
