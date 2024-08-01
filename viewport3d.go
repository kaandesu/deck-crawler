package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (gs *State) camera3DInit() {
	cam := rl.Camera{
		Position:   gs.camera.Position,
		Target:     gs.camera.Target,
		Up:         gs.camera.Up,
		Fovy:       gs.camera.Fovy,
		Projection: rl.CameraPerspective,
	}

	camOrt := rl.Camera{
		Position:   gs.camera.Position,
		Target:     gs.camera.Target,
		Up:         gs.camera.Up,
		Fovy:       gs.camera.Fovy,
		Projection: rl.CameraPerspective,
	}
	if !gs.editMode {
		rl.BeginMode3D(cam)
	} else {
		// fmt.Println("Editor deceted: CameraPerspective -> CameraOrthographic")
		rl.BeginMode3D(camOrt)
	}
}

const (
	wallH = 9
)

var test float32 = 0.0

func (gs *State) render3DViewport() {
	rl.ClearBackground(rl.SkyBlue)
	gs.camera3DInit()
	rl.DrawPlane(rl.NewVector3(0, wallH/-2, 0), rl.NewVector2(32, 32), rl.DarkBrown)

	for _, item := range ViewportState.Items {
		rl.DrawModel(item.model, item.pos, item.scale, rl.White)
	}

	rl.DrawGrid(50, 0)
	rl.EndMode3D()
}

func draw3DViewport() {
	if GameState.editMode {
		rl.BeginScissorMode(0, 0, GameScreen.width, GameScreen.height)
	} else {
		rl.BeginScissorMode(int32(GameStyle.padding), int32(GameStyle.padding), (GameScreen.width-int32(GameStyle.padding))/2, GameScreen.height/2)
	}

	GameState.render3DViewport()
	rl.EndScissorMode()
}
