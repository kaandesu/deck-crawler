package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (gs *State) camera3DInit() {
	rl.BeginMode3D(rl.Camera{
		Position:   gs.camera.Position,
		Target:     gs.camera.Target,
		Up:         gs.camera.Up,
		Fovy:       gs.camera.Fovy,
		Projection: rl.CameraPerspective,
	})
}

const (
	wallH = 9
)

func (gs *State) render3DViewport() {
	rl.ClearBackground(rl.SkyBlue)
	gs.camera3DInit()
	rl.DrawPlane(rl.NewVector3(0, wallH/-2, 0), rl.NewVector2(32, 32), rl.DarkBrown)
	rl.DrawModel(WallDoorModel, rl.NewVector3(0, -3, -10), 2.2, rl.White)

	rl.DrawModel(*WallModel, rl.NewVector3(-4, 0, 0), 2.2, rl.DarkGray)
	rl.DrawModel(*WallModel, rl.NewVector3(4, 0, 0), 2.2, rl.DarkGray)

	rl.DrawGrid(50, 0)
	rl.EndMode3D()
}
