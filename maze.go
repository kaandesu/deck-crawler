package main

import (
	"fmt"
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
	origin *Node
	matrix [][]*Node
	scale  int
}

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

func CreateMatrix(n int, scale int) Maze {
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
	return Maze{
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
			// NOTE: 4 * 2.2 = 8.8
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

var scale float32 = 2.2

func (maze *Maze) drawWalls() {
	matix := maze.matrix
	for i, row := range matix {
		for j, node := range row {
			dirs, n := node.linkNum()
			nodePos := rl.NewVector3(float32(node.X*maze.scale), 0, float32(node.Y*maze.scale))
			id := fmt.Sprintf("wall%d%d%d", i, j, n)
			switch n {
			case 1:
				offset := rl.NewVector3(0, 0, 0)
				switch dirs[0] {
				case Left:
					offset.Z = 4.4
					Scene.AddModel(Wall, id, rl.Vector3Add(nodePos, offset), rl.NewVector3(0, 0, 0), scale)
				case Right:
					offset.Z = -4.4
					Scene.AddModel(Wall, id, rl.Vector3Add(nodePos, offset), rl.NewVector3(0, 0, 0), scale)
				case Down:
					offset.X = -4.4
					Scene.AddModel(Wall, id, rl.Vector3Add(nodePos, offset), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				case Up:
					offset.X = 4.4
					Scene.AddModel(Wall, id, rl.Vector3Add(nodePos, offset), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				}
			}
		}
	}
}

func (maze *Maze) linkIncomingNodes() {
	matix := maze.matrix
	for _, row := range matix {
		for _, node := range row {
			node.linkIncoming()
		}
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
