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

type Node struct {
	Left, Right, Up, Down         *Node
	OnLeft, OnRight, OnUp, OnDown *Node
	X, Y                          int
	posX, posY                    int
	Color                         color.RGBA
}

func NewNode(x, y, scale int) *Node {
	return &Node{
		X:     x,
		Y:     y,
		posX:  x * scale,
		posY:  y * scale,
		Color: rl.Black,
	}
}

type PairNode struct {
	inBetween *Node
	nodes     []*Node
	dir       Direction // TODO: add onLeft etc dirs to the type
}

type Maze struct {
	origin    *Node
	matrix    [][]*Node
	nodePairs []PairNode
	scale     int
}

func CreateMatrix(n int, scale int) *Maze {
	grid := make([][]*Node, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]*Node, n)
		for j := 0; j < n; j++ {
			grid[i][j] = NewNode(i, j, scale)
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
	var allNodes []*Node
	for _, row := range maze.matrix {
		allNodes = append(allNodes, row...)
	}

	for _, pair := range maze.nodePairs {
		allNodes = append(allNodes, pair.inBetween)
	}

	for _, node := range allNodes {

		rl.DrawCube(rl.NewVector3(float32(node.posX), 0, float32(node.posY)), 1.5, 1.5, 1.5, node.Color)

		if node.Right != nil {
			rl.DrawLine3D(
				rl.NewVector3(float32(node.posX), 0, float32(node.posY)),
				rl.NewVector3(float32(node.Right.posX), 0, float32(node.Right.posY)),
				rl.Blue,
			)
		}

		if node.Left != nil {
			rl.DrawLine3D(
				rl.NewVector3(float32(node.posX), 0, float32(node.posY)),
				rl.NewVector3(float32(node.Left.posX), 0, float32(node.Left.posY)),
				rl.Green,
			)
		}

		if node.Up != nil {
			rl.DrawLine3D(
				rl.NewVector3(float32(node.posX), 0, float32(node.posY)),
				rl.NewVector3(float32(node.Up.posX), 0, float32(node.Up.posY)),
				rl.Blue,
			)
		}

		if node.Down != nil {
			rl.DrawLine3D(
				rl.NewVector3(float32(node.posX), 0, float32(node.posY)),
				rl.NewVector3(float32(node.Down.posX), 0, float32(node.Down.posY)),
				rl.Green,
			)
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

func (m *Maze) addInBetweenNode(pairNode PairNode) {
	pair := pairNode.nodes
	dir := pairNode.dir
	var node *Node
	switch dir {

	case Left:
		node = &Node{
			Left:  pair[1],
			Right: pair[0],
			posX:  (pair[0].posX + pair[1].posX) / 2,
			posY:  (pair[0].posY + pair[1].posY) / 2,
		}
		pair[1].Right = node
		pair[0].Left = node
		node.OnRight = pair[1]
		node.OnLeft = pair[0]

	case Right:
		node = &Node{
			Left:  pair[0],
			Right: pair[1],
			posX:  (pair[0].posX + pair[1].posX) / 2,
			posY:  (pair[0].posY + pair[1].posY) / 2,
		}
		pair[0].Right = node
		pair[1].Left = node
		node.OnRight = pair[0]
		node.OnLeft = pair[1]

	case Down:
		node = &Node{
			Up:   pair[0],
			Down: pair[1],
			posX: (pair[0].posX + pair[1].posX) / 2,
			posY: (pair[0].posY + pair[1].posY) / 2,
		}
		pair[0].Down = node
		pair[1].Up = node
		node.OnDown = pair[0]
		node.OnUp = pair[1]

	case Up:
		node = &Node{
			Up:   pair[1],
			Down: pair[0],
			posX: (pair[0].posX + pair[1].posX) / 2,
			posY: (pair[0].posY + pair[1].posY) / 2,
		}
		pair[1].Down = node
		pair[0].Up = node
		node.OnDown = pair[1]
		node.OnUp = pair[0]

	}
	node.Color = rl.White
	pairNode.inBetween = node
	m.nodePairs = append(m.nodePairs, pairNode)
}

func (m *Maze) createNodePair(node *Node) {
	dirs, _ := node.linkNum()
	for _, dir := range dirs {
		switch dir {
		case Left:
			if node.OnLeft == nil {
				m.addNodePair(node, node.Left, Left)
			}
		case Right:
			if node.OnRight == nil {
				m.addNodePair(node, node.Right, Right)
			}
		case Up:
			if node.OnUp == nil {
				m.addNodePair(node, node.Up, Up)
			}
		case Down:
			if node.OnDown == nil {
				m.addNodePair(node, node.Down, Down)
			}

		}
	}
}

func (maze *Maze) createNodePairs() {
	for _, row := range maze.matrix {
		for _, node := range row {
			maze.createNodePair(node)
		}
	}
}

func (m *Maze) addNodePair(n1, n2 *Node, dir Direction) {
	pairExists := func(pair1, pair2 []*Node) bool {
		for _, pairs := range m.nodePairs {
			pair := pairs.nodes
			if (pair[0] == pair1[0] && pair[1] == pair1[1]) || (pair[0] == pair2[0] && pair[1] == pair2[1]) {
				return true
			}
		}
		return false
	}

	pair1 := []*Node{n1, n2}
	pair2 := []*Node{n2, n1}
	if !pairExists(pair1, pair2) {
		nodePair := PairNode{
			nodes: []*Node{n1, n2},
			dir:   dir,
		}
		m.addInBetweenNode(nodePair)
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
