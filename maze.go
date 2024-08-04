package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var col color.RGBA

type Node struct {
	Left, Right, Up, Down         *Node
	OnLeft, OnRight, OnUp, OnDown *Node
	Id                            int
	X, Y                          int
	Color                         color.RGBA
	// TODO : add spawner stuct, spawner type (none,slime etc), spawnerReady bool
}

func NewNode(id, x, y int) *Node {
	return &Node{
		Id:    id,
		X:     x,
		Y:     y,
		Color: rl.Black,
	}
}

type Maze struct {
	origin    *Node
	matrix    [][]*Node
	nodePairs [][]*Node
	scale     int
}

// FIX: There are a lot of double for loops for node actions, combine them in a single function

func LinkNodes(a, b *Node, direction Direction) {
	switch direction {
	case Left:
		a.Left = b
	case Right:
		a.Right = b
	case Up:
		a.Up = b
	case Down:
		a.Down = b
	}
}

func CreateMatrix(n int, scale int) *Maze {
	grid := make([][]*Node, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]*Node, n)
		for j := 0; j < n; j++ {
			grid[i][j] = NewNode(i*n+j, i, j)
		}
	}

	// Link the nodes
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j < n-1 { // Link to the right neighbor
				LinkNodes(grid[i][j], grid[i][j+1], Right)
			}
			if j == n-1 && i < n-1 { // Link the rightmost column node to the node below
				LinkNodes(grid[i][j], grid[i+1][j], Down)
			}
		}
	}
	grid[0][n-1].Color = rl.Red
	return &Maze{
		matrix: grid,
		origin: grid[0][n-1],
		scale:  scale,
	}
}

func (maze *Maze) walkOrigin(dir Direction) {
	var node *Node
	switch dir {
	case Left:
		if maze.origin.Y == 0 {
			return
		}
		node = maze.matrix[maze.origin.X][maze.origin.Y-1]
		LinkNodes(maze.origin, node, Left)
		node.removeOriginPointer()
	case Right:
		if maze.origin.Y == len(maze.matrix[0])-1 {
			return
		}
		node = maze.matrix[maze.origin.X][maze.origin.Y+1]
		LinkNodes(maze.origin, node, Right)
		node.removeOriginPointer()
	case Up:
		if maze.origin.X == 0 {
			return
		}
		node = maze.matrix[maze.origin.X-1][maze.origin.Y]
		LinkNodes(maze.origin, node, Up)
		node.removeOriginPointer()
	case Down:
		if maze.origin.X == len(maze.matrix[0])-1 {
			return
		}
		node := maze.matrix[maze.origin.X+1][maze.origin.Y]
		LinkNodes(maze.origin, node, Down)
		node.removeOriginPointer()
	}
}

func (node *Node) removeOriginPointer() {
	if node != nil {
		maze.origin.Color = rl.Black
		if node.Right != nil {
			node.Right = nil
		}
		if node.Up != nil {
			node.Up = nil
		}
		if node.Left != nil {
			node.Left = nil
		}
		if node.Down != nil {
			node.Down = nil
		}
		node.Color = rl.Red
		maze.origin = node
	}
}

func (maze *Maze) draw() {
	for _, row := range maze.matrix {
		for _, node := range row {
			col = node.Color
			rl.DrawCube(rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale)), 1.5, 1.5, 1.5, col)

			if node.Right != nil {
				rl.DrawLine3D(
					rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale)),
					rl.NewVector3(float32(node.Right.X*maze.scale), 0, float32(node.Right.Y*maze.scale)),
					rl.Blue,
				)
			}

			if node.Left != nil {
				rl.DrawLine3D(
					rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale)),
					rl.NewVector3(float32(node.Left.X*maze.scale), 0, float32(node.Left.Y*maze.scale)),
					rl.Green,
				)
			}

			if node.Up != nil {
				rl.DrawLine3D(
					rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale)),
					rl.NewVector3(float32(node.Up.X*maze.scale), 0, float32(node.Up.Y*maze.scale)),
					rl.Blue,
				)
			}

			if node.Down != nil {
				rl.DrawLine3D(
					rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale)),
					rl.NewVector3(float32(node.Down.X*maze.scale), 0, float32(node.Down.Y*maze.scale)),
					rl.Green,
				)
			}

		}
	}
}

var (
	scale    float32 = 2.2
	baseSize float32 = 4
	wallSize float32 = scale * baseSize
)

func (maze *Maze) drawWalls() {
	for _, row := range maze.matrix {
		for _, node := range row {
			node.linkIncoming()
			maze.createNodePair(node)
			dirs, n := node.linkNum()
			nodePos := rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale))
			switch n {
			case 1:
				switch dirs[0] {
				case Left:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, wallSize/2)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-(wallSize/2+0.5), 0, 0.5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3((wallSize/2+0.5), 0, 0.5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				case Right:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, -wallSize/2)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-(wallSize/2+0.5), 0, 0.5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3((wallSize/2+0.5), 0, 0.5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				case Down:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-wallSize/2, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0.5, 0, -(wallSize/2+0.5))), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0.5, 0, wallSize/2+0.5)), rl.NewVector3(0, 0, 0), scale)
				case Up:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(wallSize/2, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-0.5, 0, -(wallSize/2+0.5))), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-0.5, 0, wallSize/2+0.5)), rl.NewVector3(0, 0, 0), scale)
				}
			}
		}
	}
}

// WARN: will use this to create the in-between nodes
// will use them to create more walls and walking points
func (m *Maze) createNodePair(node *Node) {
	dirs, _ := node.linkNum()
	for _, dir := range dirs {
		switch dir {
		case Left:
			if node.OnLeft != nil {
				m.addNodePair(node, node.OnLeft)
			} else {
				m.addNodePair(node, node.Left)
			}
		case Right:
			if node.OnRight != nil {
				m.addNodePair(node, node.OnRight)
			} else {
				m.addNodePair(node, node.Right)
			}
		case Up:
			if node.OnUp != nil {
				m.addNodePair(node, node.OnUp)
			} else {
				m.addNodePair(node, node.Up)
			}
		case Down:
			if node.OnDown != nil {
				m.addNodePair(node, node.OnDown)
			} else {
				m.addNodePair(node, node.Down)
			}

		}
	}
}

func (m *Maze) addNodePair(n1, n2 *Node) {
	pairExists := func(pair1, pair2 []*Node) bool {
		for _, pair := range m.nodePairs {
			if (pair[0] == pair1[0] && pair[1] == pair1[1]) || (pair[0] == pair2[0] && pair[1] == pair2[1]) {
				return true
			}
		}
		return false
	}

	pair1 := []*Node{n1, n2}
	pair2 := []*Node{n2, n1}
	if !pairExists(pair1, pair2) {
		m.nodePairs = append(m.nodePairs, pair1)
	}
}

func (node *Node) linkIncoming() {
	if node.X != 0 {
		n := maze.matrix[node.X-1][node.Y]
		if n.Down == node {
			node.OnUp = n
		}
	}
	if node.X != len(maze.matrix)-1 {
		n := maze.matrix[node.X+1][node.Y]
		if n.Up == node {
			node.OnDown = n
		}
	}

	if node.Y != len(maze.matrix)-1 {
		n := maze.matrix[node.X][node.Y+1]
		if n.Left == node {
			node.OnRight = n
		}
	}

	if node.Y != 0 {
		n := maze.matrix[node.X][node.Y-1]
		if n.Right == node {
			node.OnLeft = n
		}
	}
}

func (node *Node) linkNum() ([]Direction, int) {
	var (
		dirs  []Direction
		count = 0
	)

	if node.Left != nil {
		dirs = append(dirs, Left)
		count++
	} else if node.OnLeft != nil {
		dirs = append(dirs, Left)
		count++
	}

	if node.Right != nil {
		dirs = append(dirs, Right)
		count++
	} else if node.OnRight != nil {
		dirs = append(dirs, Right)
		count++
	}

	if node.Down != nil {
		dirs = append(dirs, Down)
		count++
	} else if node.OnDown != nil {
		dirs = append(dirs, Down)
		count++
	}

	if node.Up != nil {
		dirs = append(dirs, Up)
		count++
	} else if node.OnUp != nil {
		dirs = append(dirs, Up)
		count++
	}
	switch count {
	case 0:
		node.Color = rl.DarkGray
	case 1:
		node.Color = rl.Red
	case 2:
		node.Color = rl.Yellow
	case 3:
		node.Color = rl.Green
	case 4:
		node.Color = rl.Violet
	}

	return dirs, count
}