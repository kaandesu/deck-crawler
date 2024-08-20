package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Screen struct {
	title  string
	width  int32
	height int32
	fps    int32
}

type Style struct {
	bg        color.RGBA
	primary   color.RGBA
	accent    color.RGBA
	padding   float32
	roundness float32
}

type State struct {
	camera             *Camera
	running            bool
	editMode           bool
	editFull           bool
	editFocusedItemUid string
	camMode            rl.CameraMode
	lookDir            Direction
	currentNode        *Node
}

type SceneItem struct {
	uid       string
	model     rl.Model
	pos       rl.Vector3
	rot       rl.Vector3
	scale     float32
	hidden    bool
	highlight bool
}

type Scene3D struct {
	Items   map[string]*SceneItem
	Enemies Enemies
}

type Enemies = map[EnemyType]Enemy
