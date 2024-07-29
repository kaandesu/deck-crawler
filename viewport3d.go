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

func (gs *State) render3DViewport() {
	rl.ClearBackground(rl.RayWhite)
	gs.camera3DInit()
	rl.DrawCube(rl.NewVector3(0, 0, 0), 2, 2, 2, rl.Red)
	rl.DrawCube(rl.NewVector3(4, 0, 0), 2, 2, 2, rl.Green)
	rl.DrawCube(rl.NewVector3(-4, 0, 0), 2, 2, 2, rl.Blue)
	rl.DrawCube(rl.NewVector3(0, 0, 4), 2, 2, 2, rl.Yellow)
	rl.DrawCube(rl.NewVector3(0, 0, -4), 2, 2, 2, rl.Purple)
	rl.DrawGrid(10, 1.0)
	rl.EndMode3D()
}
