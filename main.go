package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var GameStyle = &Style{
	bg:        rl.NewColor(50, 50, 50, 255),
	primary:   rl.NewColor(200, 200, 200, 255),
	accent:    rl.NewColor(200, 200, 0, 255),
	padding:   20,
	roundness: 0.1,
}

var GameScreen = &Screen{
	width:  800,
	height: 680,
	title:  "Deck Crawler",
	fps:    60,
}

var GameState = &State{
	running: true,
	camera:  NewCamera(),
}

var ViewportState = &Scene3D{
	Items: make(map[string]*SceneItem),
}

var button bool

func drawScene() {
	rl.DrawText("Deck Crawler!", GameScreen.width/2+int32(GameStyle.padding)*2, int32(GameStyle.padding), 45, GameStyle.accent)
	draw3DViewport()
}

func main() {
	server := NewServer(":3000")
	go server.Start()
	defer quit()
	setup()

	for GameState.running {
		input()
		update()
		render()
	}
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
	LoadModels()
}
