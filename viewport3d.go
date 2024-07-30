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
	rl.DrawCube(rl.NewVector3(4, 0, 0), 0.2, wallH, 20, rl.DarkGray)
	rl.DrawCube(rl.NewVector3(-4, 0, 0), 0.2, wallH, 20, rl.DarkPurple)
	rl.DrawCube(rl.NewVector3(0, 0, -10), 9, wallH, 0.2, rl.DarkBlue)
	rl.DrawGrid(50, 0)
	rl.EndMode3D()
}
