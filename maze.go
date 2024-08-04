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
		b.OnRight = a
	case Right:
		a.Right = b
		b.OnLeft = a
	case Up:
		a.Up = b
		b.OnDown = a
	case Down:
		a.Down = b
		b.OnUp = a
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
		maze.origin = node
		maze.origin.Color = rl.Red
		if maze.origin.Right != nil {
			maze.origin.Right = nil
		}
		if maze.origin.Up != nil {
			maze.origin.Up = nil
		}
		if maze.origin.Left != nil {
			maze.origin.Left = nil
		}
		if maze.origin.Down != nil {
			maze.origin.Down = nil
		}
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
