package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	padding   = 10
	roundness = 0.1
)

type Screen struct {
	title  string
	width  int32
	height int32
	fps    int32
}

type Style struct {
	bg      color.RGBA
	primary color.RGBA
	accent  color.RGBA
	padding float32
}

var GameStyle = &Style{
	bg:      rl.NewColor(100, 100, 100, 255),
	primary: rl.NewColor(200, 200, 200, 255),
	accent:  rl.NewColor(200, 200, 0, 255),
	padding: 20,
}

var GameScreen = &Screen{
	width:  800,
	height: 680,
	title:  "Deck Crawler",
	fps:    60,
}

type State struct {
	camera  *Camera
	running bool
}

var GameState = &State{
	running: true,
	camera:  NewCamera(),
}

var button bool

func drawScene() {
	// rl.DrawRectangleRoundedLines(rl.NewRectangle(float32(GameScreen.width)/2, GameStyle.padding, float32(GameScreen.width)/2-GameStyle.padding, float32(GameScreen.width)/2-GameStyle.padding), 0.1, 0, 1, GameStyle.accent)
	rl.BeginScissorMode(0, 0, GameScreen.width/2, GameScreen.height/2)
	GameState.render3DViewport()
	rl.EndScissorMode()
}

func input() {
	rl.UpdateCamera(GameState.camera, rl.CameraFirstPerson)
}

func update() {
	GameState.running = !rl.WindowShouldClose()
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(GameStyle.bg)
	drawScene()
	rl.EndDrawing()
}

func quit() {
	rl.CloseWindow()
}

func setup() {
	rl.InitWindow(GameScreen.width, GameScreen.height, GameScreen.title)
	rl.SetExitKey(0)
	rl.SetTargetFPS(GameScreen.fps)
}

func main() {
	defer quit()
	setup()
	for GameState.running {
		input()
		update()
		render()
	}
}
