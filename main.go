package main

// TODO:
// - block movement forward / backward if next has an enemy
// - spawn enemy after a minute

import (
	"errors"
	"flag"
	"fmt"
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
		editMode:           false,
		editFull:           true,
		editFocusedItemUid: "",
		camMode:            rl.CameraFree,
		lookDir:            Right,
		currentNode:        &Node{},
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
	dir := ""
	switch GameState.lookDir {
	case Left:
		dir = "Left"
	case Right:
		dir = "Right"
	case Up:
		dir = "Up"
	case Down:
		dir = "Down"
	}
	rl.DrawText(fmt.Sprintf("[ %s -  %d ] --- %+v", dir, GameState.lookDir, GameState.camera.Position), 100, GameScreen.height-50, 25, GameStyle.accent)
}

var shader rl.Shader

func main() {
	flag.Parse()
	if *enableEditorServer {
		server := NewServer("127.0.0.1:3000")
		go server.Start()
	}
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	setup()

	for GameState.running {
		input()
		update()
		render()
	}

	defer quit()
	defer rl.UnloadRenderTexture(SceneRenderTexture)
	defer rl.UnloadShader(renderShader)
	defer rl.UnloadShader(shader)
	defer unloadEnemyTextures()
	defer rl.UnloadModel(ground)
	defer rl.UnloadModel(testCube)
}

func unloadEnemyTextures() {
	for _, enemy := range Scene.Enemies {
		rl.UnloadTexture(enemy.Texture)
	}
}

var (
	inputBlocked    bool
	turningLeft     bool
	turningRight    bool
	movingForward   bool
	movingBackward  bool
	currentRotation float32
)

func input() {
	if rl.IsKeyPressed(rl.KeyY) {
		lights[0].enabled *= -1
	}
	if rl.IsKeyPressed(rl.KeyR) {
		lights[1].enabled *= -1
	}
	if rl.IsKeyPressed(rl.KeyG) {
		lights[2].enabled *= -1
	}
	if rl.IsKeyPressed(rl.KeyB) {
		lights[3].enabled *= -1
	}

	for i := 0; i < len(lights); i++ {
		lights[i].UpdateValues()
	}

	if rl.IsKeyDown(rl.KeyA) {
		if inputBlocked || movingForward || movingBackward {
			return
		}
		GameState.lookDir = Direction((GameState.lookDir + 1) % 4)
		turningLeft = true
		currentRotation = 0.0
		blockInputs()
	}

	if rl.IsKeyDown(rl.KeyD) {
		if inputBlocked || movingForward || movingBackward {
			return
		}

		switch GameState.lookDir {
		case Left:
			GameState.lookDir = Down
		case Right:
			GameState.lookDir = Up
		case Down:
			GameState.lookDir = Right
		case Up:
			GameState.lookDir = Left
		}

		turningRight = true
		currentRotation = 0.0
		blockInputs()
	}

	if rl.IsKeyDown(rl.KeyW) {
		if movingToNode || turningLeft || turningRight {
			return
		}

		dirs := GameState.currentNode.linkNum()
		switch GameState.lookDir {
		case Left:
			movingForward = includeDir(dirs, []Direction{Left})
			if movingForward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Left != nil {
					GameState.currentNode = GameState.currentNode.Left
				} else {
					GameState.currentNode = GameState.currentNode.OnLeft
				}
			}

		case Right:
			movingForward = includeDir(dirs, []Direction{Right})
			if movingForward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Right != nil {
					GameState.currentNode = GameState.currentNode.Right
				} else {
					GameState.currentNode = GameState.currentNode.OnRight
				}
			}

		case Up:
			movingForward = includeDir(dirs, []Direction{Up})
			if movingForward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Up != nil {
					GameState.currentNode = GameState.currentNode.Up
				} else {
					GameState.currentNode = GameState.currentNode.OnUp
				}
			}

		case Down:
			movingForward = includeDir(dirs, []Direction{Down})
			if movingForward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Down != nil {
					GameState.currentNode = GameState.currentNode.Down
				} else {
					GameState.currentNode = GameState.currentNode.OnDown
				}
			}

		}

		if movingForward {
			targetPos = rl.NewVector3(GameState.currentNode.posX, GameState.camera.Position.Y, GameState.currentNode.posY)
			movingToNode = true
		}

	}

	if rl.IsKeyDown(rl.KeyS) {
		if movingToNode || turningLeft || turningRight {
			return
		}

		dirs := GameState.currentNode.linkNum()
		switch GameState.lookDir {
		case Right:
			movingBackward = includeDir(dirs, []Direction{Left})
			if movingBackward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Left != nil {
					GameState.currentNode = GameState.currentNode.Left
				} else {
					GameState.currentNode = GameState.currentNode.OnLeft
				}
			}

		case Left:
			movingBackward = includeDir(dirs, []Direction{Right})
			if movingBackward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Right != nil {
					GameState.currentNode = GameState.currentNode.Right
				} else {
					GameState.currentNode = GameState.currentNode.OnRight
				}
			}

		case Down:
			movingBackward = includeDir(dirs, []Direction{Up})
			if movingBackward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Up != nil {
					GameState.currentNode = GameState.currentNode.Up
				} else {
					GameState.currentNode = GameState.currentNode.OnUp
				}
			}

		case Up:
			movingBackward = includeDir(dirs, []Direction{Down})
			if movingBackward {
				GameState.currentNode.Color = rl.Brown
				if GameState.currentNode.Down != nil {
					GameState.currentNode = GameState.currentNode.Down
				} else {
					GameState.currentNode = GameState.currentNode.OnDown
				}
			}

		}

		if movingBackward {
			targetPos = rl.NewVector3(GameState.currentNode.posX, GameState.camera.Position.Y, GameState.currentNode.posY)
			movingToNode = true
		}

	}
}

func blockInputs() {
	if !inputBlocked {
		inputBlocked = true
		go func() {
			time.Sleep(time.Millisecond * 500)
			turningLeft = false
			turningRight = false
			inputBlocked = false
			movingBackward = false
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

	UpdateCameraCustom(GameState.camera)

	rl.DrawFPS(GameScreen.width-90, GameScreen.height-30)
	rl.EndDrawing()
}

func quit() {
	rl.CloseWindow()
}

var lights []Light

func setup() {
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(GameScreen.width, GameScreen.height, GameScreen.title)
	rl.SetExitKey(0)
	rl.SetTargetFPS(GameScreen.fps)

	ground = rl.LoadModelFromMesh(rl.GenMeshPlane(10, 10, 3, 3))
	shader = rl.LoadShader("./res/shaders/glsl330/lighting.vs", "./res/shaders/glsl330/lighting.fs")
	*shader.Locs = rl.GetShaderLocation(shader, "viewPos")
	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	shaderValue := []float32{0.1, 0.1, 0.1, 1.0}
	rl.SetShaderValue(shader, ambientLoc, shaderValue, rl.ShaderUniformVec4)
	ground.Materials.Shader = shader

	testCube = rl.LoadModelFromMesh(rl.GenMeshCube(3, 3, 3))
	testCube.Materials.Shader = shader

	lights = make([]Light, 4)
	lights[0] = NewLight(LightTypePoint, rl.NewVector3(0, 3, 30), rl.NewVector3(0, 0, 0), rl.Yellow, shader)
	lights[1] = NewLight(LightTypePoint, rl.NewVector3(0, 3, 40), rl.NewVector3(0, 0, 0), rl.Red, shader)
	lights[2] = NewLight(LightTypePoint, rl.NewVector3(0, 3, 60), rl.NewVector3(0, 0, 0), rl.Green, shader)
	lights[3] = NewLight(LightTypePoint, rl.NewVector3(0, 3, 50), rl.NewVector3(0, 0, 0), rl.Blue, shader)

	Scene.Enemies = map[EnemyType]Enemy{
		Slime:    DefineEnemy(Slime, []float32{1, 2}, 5, 0, "./res/imgs/billboard.png"),
		NotSlime: DefineEnemy(NotSlime, []float32{1, 2}, 5, 0, "./res/imgs/billboard.png"),
	}

	// renderShader = rl.LoadShader("./res/shaders/glsl330/base.vs", "./res/shaders/glsl330/cross_stitching.fs")
	// renderShader = rl.LoadShader("./res/shaders/glsl330/base.vs", "./res/shaders/glsl330/pixelizer.fs")
	renderShader = rl.LoadShader("./res/shaders/glsl330/base.vs", "./res/shaders/glsl330/base.fs")
	maze = CreateMatrix(5, 17.6)
	for range len(maze.matrix) * len(maze.matrix) * 12 {
		maze.walkOrigin(Direction(rand.Intn(4)))
	}
	maze.drawWalls()
	maze.createNodePairs()
	GameState.camera = NewCamera()
	maze.setAllNodes()
	GameState.currentNode = maze.matrix[0][0]

	maze.drawInBetweenWallPairs()

	// TODO: change this with assert
	if len(maze.nodePairs) != len(maze.matrix)*len(maze.matrix)-1 {
		panic(errors.New("pair num is not correct"))
	}

	if *enableFullScreen {
		SceneRenderTexture = rl.LoadRenderTexture(GameScreen.width*6/7, GameScreen.height*6/7)
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
