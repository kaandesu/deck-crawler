package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var cameraInit bool = false
var (
	targetOffset   = rl.NewVector3(2, 9, -0.85)
	positionOffset = rl.NewVector3(3.6, 45, 0)
)

func (gs *State) camera3DInit() {
	gs.camera.Up = rl.NewVector3(0, 1, 0)
	// rl.UpdateCamera(gs.camera, gs.camMode)
	if gs.camMode == rl.CameraThirdPerson {
		if cameraInit {
			targetOffset = rl.NewVector3(0, 0, 0)
			positionOffset = rl.NewVector3(0, 0, 0)
		}
		gs.camera.Position = rl.Vector3Add(gs.camera.Position, positionOffset)
		gs.camera.Target = rl.Vector3Add(gs.camera.Target, targetOffset)
		cameraInit = true
	}
	// fmt.Printf("%+v \n", gs.camera)
}

const (
	wallH = 9
)

var test float32 = 0.0

func (gs *State) render3DViewport() {
	rl.ClearBackground(rl.SkyBlue)
	gs.camera3DInit()
	rl.BeginMode3D(*gs.camera)

	rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(32, 32), rl.DarkBrown)

	for _, item := range ViewportState.Items {
		if item.hidden {
			continue
		}
		rl.DrawModel(item.model, item.pos, item.scale, rl.White)
	}

	rl.DrawGrid(50, 0)
	rl.EndMode3D()
}

func draw3DViewport() {
	if GameState.editMode && GameState.editFull {
		rl.BeginScissorMode(0, 0, GameScreen.width, GameScreen.height)
	} else {
		rl.BeginScissorMode(int32(GameStyle.padding), int32(GameStyle.padding), (GameScreen.width-int32(GameStyle.padding))/2, GameScreen.height/2)
	}

	GameState.render3DViewport()
	rl.EndScissorMode()
}
