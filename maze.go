package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up

	OnLeft
	OnRight
	OnUp
	OnDown
)

type Node struct {
	Left, Right, Up, Down         *Node
	OnLeft, OnRight, OnUp, OnDown *Node
	X, Y                          int
	posX, posY                    float32
	Color                         color.RGBA
}

func NewNode(x, y int, scale float32) *Node {
	return &Node{
		X:     x,
		Y:     y,
		posX:  float32(x) * scale,
		posY:  float32(y) * scale,
		Color: rl.Black,
	}
}

type PairNode struct {
	inBetween *Node
	nodes     []*Node
	dir       Direction
}

type Maze struct {
	origin    *Node
	matrix    [][]*Node
	nodePairs []PairNode
	allNodes  []*Node
	scale     float32
}

func CreateMatrix(n int, scale float32) *Maze {
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
	for _, node := range maze.allNodes {

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
			dirs := node.linkNum()
			nodePos := rl.NewVector3(float32(node.X)*maze.scale, 0, float32(node.Y)*maze.scale)
			switch len(dirs) {
			case 1:
				switch dirs[0] {
				case Left:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, wallSize/2)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-(wallSize/2+1.1), 0, -2.45)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3((wallSize/2+1.1), 0, -2.45)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, 4.15)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, 4.15)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				case Right:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, -wallSize/2)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-(wallSize/2+1.1), 0, 2.45)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3((wallSize/2+1.1), 0, 2.45)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, -4.15)), rl.NewVector3(0, 0*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, -4.15)), rl.NewVector3(0, -90*rl.Deg2rad, 0), scale)
				case Down:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-wallSize/2, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(2.45, 0, -(wallSize/2+0.55))), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(2.45, 0, wallSize/2+0.55)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-4.15, 0, 5)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-4.15, 0, -5)), rl.NewVector3(0, -90*rl.Deg2rad, 0), scale)
				case Up:
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(wallSize/2, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-2.45, 0, -(wallSize/2+0.55))), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-2.45, 0, wallSize/2+0.55)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(4.15, 0, -5)), rl.NewVector3(0, 0*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(4.15, 0, 5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				}
			case 2:
				if includeDir(dirs, []Direction{Left, Down}) { // NOTE: yellow corners
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, -5)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, wallSize/2+0.6)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-wallSize/2-1.1, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, 5)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
				} else if includeDir(dirs, []Direction{Left, Up}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, -5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, wallSize/2+0.6)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(wallSize/2+1.1, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, 5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				} else if includeDir(dirs, []Direction{Right, Down}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, 5)), rl.NewVector3(0, 270*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-wallSize/2-1.1, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, -wallSize/2-0.6)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, -5)), rl.NewVector3(0, 270*rl.Deg2rad, 0), scale)
				} else if includeDir(dirs, []Direction{Right, Up}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, 5)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, -wallSize/2-0.6)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(wallSize/2+1.1, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, -5)), rl.NewVector3(0, 0*rl.Deg2rad, 0), scale)
				} else if includeDir(dirs, []Direction{Right, Left}) { // NOTE: yellow non-corners
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(3.3-wallSize, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3((wallSize)-3.3, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				} else if includeDir(dirs, []Direction{Up, Down}) {
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, 3.8-wallSize)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, (wallSize)-3.8)), rl.NewVector3(0, 0, 0), scale)
				}
			case 3:
				if !includeDir(dirs, []Direction{Left}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, 5)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, 5)), rl.NewVector3(0, 270*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, -wallSize/2-0.6)), rl.NewVector3(0, 0, 0), scale)
				} else if !includeDir(dirs, []Direction{Right}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, -5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, -5)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, wallSize/2+0.6)), rl.NewVector3(0, 0, 0), scale)
				} else if !includeDir(dirs, []Direction{Down}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, 5)), rl.NewVector3(0, 0, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, -5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(wallSize/2+1.1, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				} else if !includeDir(dirs, []Direction{Up}) {
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, -5)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
					Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, 5)), rl.NewVector3(0, 270*rl.Deg2rad, 0), scale)
					Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(-wallSize/2-1.1, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				}

			case 4:
				Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, 5)), rl.NewVector3(0, 0, 0), scale)
				Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(-5.5, 0, -5)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
				Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, -5)), rl.NewVector3(0, 180*rl.Deg2rad, 0), scale)
				Scene.AddModel(WallCornerSmall, rl.Vector3Add(nodePos, rl.NewVector3(5.5, 0, 5)), rl.NewVector3(0, 270*rl.Deg2rad, 0), scale)
			}
		}
	}
}

func includeDir(dirs []Direction, directions []Direction) bool {
	dirMap := make(map[Direction]bool)
	for _, dir := range dirs {
		dirMap[dir] = true
	}

	for _, direction := range directions {
		switch direction {
		case Left:
			if !dirMap[Left] && !dirMap[OnLeft] {
				return false
			}
		case Right:
			if !dirMap[Right] && !dirMap[OnRight] {
				return false
			}
		case Up:
			if !dirMap[Up] && !dirMap[OnUp] {
				return false
			}
		case Down:
			if !dirMap[Down] && !dirMap[OnDown] {
				return false
			}
		default:
			if !dirMap[direction] {
				return false
			}
		}
	}

	return true
}

func (node *Node) drawInBetweenWalls() {
	nodePos := rl.NewVector3(float32(node.posX), 0, float32(node.posY))
	_ = nodePos
	if node.Left != nil || node.Right != nil {
		node.Color = rl.Pink
		Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(3.3-wallSize, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
		Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3((wallSize)-3.3, 0, 0)), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
	}
	if node.Up != nil || node.Down != nil {
		node.Color = rl.Beige
		Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, 3.8-wallSize)), rl.NewVector3(0, 0, 0), scale)
		Scene.AddModel(Wall, rl.Vector3Add(nodePos, rl.NewVector3(0, 0, (wallSize)-3.8)), rl.NewVector3(0, 0, 0), scale)
	}
}

func (m *Maze) addInBetweenNode(pairNode PairNode) {
	pair := pairNode.nodes
	dir := pairNode.dir
	var node *Node
	switch dir {

	case Left:
		node = &Node{
			Left:    pair[1],
			OnRight: pair[0],
			posX:    (pair[0].posX + pair[1].posX) / 2,
			posY:    (pair[0].posY + pair[1].posY) / 2,
		}
		pair[0].Left = node
		pair[1].OnRight = node

	case Right:
		node = &Node{
			OnLeft: pair[0],
			Right:  pair[1],
			posX:   (pair[0].posX + pair[1].posX) / 2,
			posY:   (pair[0].posY + pair[1].posY) / 2,
		}
		pair[0].Right = node
		pair[0].Color = rl.Pink
		pair[1].OnLeft = node

	case Down:
		node = &Node{
			OnUp: pair[0],
			Down: pair[1],
			posX: (pair[0].posX + pair[1].posX) / 2,
			posY: (pair[0].posY + pair[1].posY) / 2,
		}
		pair[0].Down = node
		pair[1].OnUp = node

	case Up:
		node = &Node{
			Up:     pair[1],
			OnDown: pair[0],
			posX:   (pair[0].posX + pair[1].posX) / 2,
			posY:   (pair[0].posY + pair[1].posY) / 2,
		}

		pair[0].Up = node
		pair[1].OnDown = node
	}

	node.Color = rl.White
	pairNode.inBetween = node
	m.nodePairs = append(m.nodePairs, pairNode)
}

func (m *Maze) createNodePair(node *Node) {
	dirs := node.linkNum()
	for _, dir := range dirs {
		switch dir {
		case Left:
			m.addNodePair(node, node.Left, Left)
		case Right:
			m.addNodePair(node, node.Right, Right)
		case Up:
			m.addNodePair(node, node.Up, Up)
		case Down:
			m.addNodePair(node, node.Down, Down)

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

func (maze *Maze) drawInBetweenWallPairs() {
	for _, node := range maze.nodePairs {
		node.inBetween.drawInBetweenWalls()
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

func (node *Node) linkNum() []Direction {
	var (
		dirs  []Direction
		count = 0
	)

	if node.Left != nil {
		dirs = append(dirs, Left)
		count++
	}

	if node.OnLeft != nil {
		dirs = append(dirs, OnLeft)
		count++
	}

	if node.Right != nil {
		dirs = append(dirs, Right)
		count++
	}

	if node.OnRight != nil {
		dirs = append(dirs, OnRight)
		count++
	}

	if node.Down != nil {
		dirs = append(dirs, Down)
		count++
	}

	if node.OnDown != nil {
		dirs = append(dirs, OnDown)
		count++
	}

	if node.Up != nil {
		dirs = append(dirs, Up)
		count++
	}

	if node.OnUp != nil {
		dirs = append(dirs, OnUp)
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

	return dirs
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
