package main

import (
	"flag"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/rand"
)

var (
	GameStyle = &Style{
		bg:        rl.NewColor(50, 50, 50, 255),
		primary:   rl.NewColor(200, 200, 200, 255),
		accent:    rl.NewColor(200, 200, 0, 255),
		padding:   20,
		roundness: 0.1,
	}
	GameScreen = &Screen{
		width:  800,
		height: 680,
		title:  "Deck Crawler",
		fps:    60,
	}
	GameState = &State{
		running:            true,
		camera:             NewCamera(),
		editMode:           false,
		editFull:           true,
		editFocusedItemUid: "",
		camMode:            rl.CameraThirdPerson,
	}
	Scene = &Scene3D{
		Items: make(map[string]*SceneItem),
	}
	enableEditorServer = flag.Bool("server", false, "enable editor server")
	enableFullScreen   = flag.Bool("full", false, "enable full screen for editing")
	SceneRenderTexture rl.RenderTexture2D
	SceneRenderRect    rl.Rectangle
	maze               Maze
)

var renderShader rl.Shader

func drawScene() {
	rl.DrawText("Deck Crawler!", GameScreen.width/2+int32(GameStyle.padding)*2, int32(GameStyle.padding), 45, GameStyle.accent)
}

func main() {
	flag.Parse()
	if *enableEditorServer {
		server := NewServer("127.0.0.1:3000")
		go server.Start()
	}

	setup()
	for GameState.running {
		input()
		update()
		render()
	}

	defer quit()
	defer rl.UnloadRenderTexture(SceneRenderTexture)
	defer rl.UnloadShader(renderShader)
}

func input() {
}

func update() {
	GameState.running = !rl.WindowShouldClose()
}

func render() {
	rl.BeginTextureMode(SceneRenderTexture)
	Scene.render()
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(GameStyle.bg)
	drawScene()

	rl.BeginShaderMode(renderShader)
	rl.DrawTextureRec(SceneRenderTexture.Texture, SceneRenderRect, rl.NewVector2(GameStyle.padding, GameStyle.padding), rl.White)
	rl.EndShaderMode()

	rl.UpdateCamera(GameState.camera, GameState.camMode)

	rl.DrawFPS(GameScreen.width-90, GameScreen.height-30)
	rl.EndDrawing()
}

func quit() {
	rl.CloseWindow()
}

func setup() {
	rl.InitWindow(GameScreen.width, GameScreen.height, GameScreen.title)
	rl.SetExitKey(0)
	rl.SetTargetFPS(GameScreen.fps)
	// renderShader = rl.LoadShader("./res/shaders/glsl330/base.vs", "./res/shaders/glsl330/cross_stitching.fs")
	renderShader = rl.LoadShader("./res/shaders/glsl330/base.vs", "./res/shaders/glsl330/base.fs")
	maze = CreateMatrix(8, 24)
	for range len(maze.matrix) * len(maze.matrix) * 11 {
		maze.walkOrigin(Direction(rand.Intn(4)))
	}
	maze.linkIncomingNodes()
	maze.drawWalls()

	if *enableFullScreen {
		SceneRenderTexture = rl.LoadRenderTexture(GameScreen.width, GameScreen.height)
	} else {
		SceneRenderTexture = rl.LoadRenderTexture(GameScreen.width/2, GameScreen.height/2+int32(GameStyle.padding)*2)
	}
	SceneRenderRect = rl.NewRectangle(0, 0, float32(SceneRenderTexture.Texture.Width), -float32(SceneRenderTexture.Texture.Height))
	LoadModels()
}
