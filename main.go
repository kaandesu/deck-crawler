package main

import (
	"errors"
	"flag"
	"time"

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
	maze               *Maze
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

var (
	inputBlocked    bool
	turningLeft     bool
	turningRight    bool
	targetRotation  float32
	currentRotation float32
	rotationSpeed   float32 = 90.0 // Degrees per second
)

func input() {
	if rl.IsKeyDown(rl.KeyA) {
		if inputBlocked {
			return
		}
		turningLeft = true
		targetRotation = 90.0 // Rotate 90 degrees
		currentRotation = 0.0
		blockInputs()
	}

	if rl.IsKeyDown(rl.KeyD) {
		if inputBlocked {
			return
		}
		turningRight = true
		targetRotation = -90.0 // Rotate -90 degrees
		currentRotation = 0.0
		blockInputs()
	}
}

func blockInputs() {
	if !inputBlocked {
		inputBlocked = true
		go func() {
			time.Sleep(time.Second * 1)
			turningLeft = false
			turningRight = false
			inputBlocked = false
		}()
	}
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

	// rl.UpdateCamera(GameState.camera, GameState.camMode)
	//
	UpdateCameraCustom(GameState.camera)

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
	maze = CreateMatrix(5, 17.6)
	for range len(maze.matrix) * len(maze.matrix) * 11 {
		maze.walkOrigin(Direction(rand.Intn(4)))
	}
	maze.drawWalls()
	maze.createNodePairs()
	maze.setAllNodes()
	maze.drawInBetweenWallPairs()

	if len(maze.nodePairs) != len(maze.matrix)*len(maze.matrix)-1 {
		panic(errors.New("pair num is not correct"))
	}

	if *enableFullScreen {
		SceneRenderTexture = rl.LoadRenderTexture(GameScreen.width, GameScreen.height)
	} else {
		SceneRenderTexture = rl.LoadRenderTexture(GameScreen.width/2, GameScreen.height/2+int32(GameStyle.padding)*2)
	}
	SceneRenderRect = rl.NewRectangle(0, 0, float32(SceneRenderTexture.Texture.Width), -float32(SceneRenderTexture.Texture.Height))
}

func (m *Maze) setAllNodes() {
	for _, row := range maze.matrix {
		m.allNodes = append(m.allNodes, row...)
	}

	for _, pair := range maze.nodePairs {
		m.allNodes = append(m.allNodes, pair.inBetween)
	}
}
