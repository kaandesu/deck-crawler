package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera = rl.Camera3D

func NewCamera() *Camera {
	return &Camera{
		Position:   rl.NewVector3(15.0, 5.0, 15.0),
		Target:     rl.NewVector3(0.0, 0.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       50.0,
		Projection: rl.CameraPerspective,
	}
}

func MoveCam(cam *Camera, direction rl.Vector3) {
	cam.Position = rl.Vector3Add(cam.Position, direction)
	// cam.Target = rl.Vector3Add(cam.Target, direction)
}

func UpdateCameraCustom(camera *rl.Camera) {
	var (
		moveSpeed   float32 = 0.2
		rotateSpeed float32 = 2.0
	)

	forwardDir := rl.Vector3Subtract(camera.Target, camera.Position)
	forwardDir.Y = 0
	forwardDir = rl.Vector3Normalize(forwardDir)

	rightDir := rl.Vector3CrossProduct(camera.Up, forwardDir)
	rightDir = rl.Vector3Normalize(rightDir)

	if rl.IsKeyDown(rl.KeyW) {
		camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(forwardDir, moveSpeed))
		camera.Target = rl.Vector3Add(camera.Target, rl.Vector3Scale(forwardDir, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		camera.Position = rl.Vector3Subtract(camera.Position, rl.Vector3Scale(forwardDir, moveSpeed))
		camera.Target = rl.Vector3Subtract(camera.Target, rl.Vector3Scale(forwardDir, moveSpeed))
	}

	if rl.IsKeyDown(rl.KeyD) {
		camera.Position = rl.Vector3Subtract(camera.Position, rl.Vector3Scale(rightDir, moveSpeed))
		camera.Target = rl.Vector3Subtract(camera.Target, rl.Vector3Scale(rightDir, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyA) {
		camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(rightDir, moveSpeed))
		camera.Target = rl.Vector3Add(camera.Target, rl.Vector3Scale(rightDir, moveSpeed))
	}

	if rl.IsKeyDown(rl.KeyQ) {
		rotateAngle := rotateSpeed * rl.GetFrameTime()
		camera.Target = rotateAround(camera.Target, camera.Position, camera.Up, rotateAngle)
	}
	if rl.IsKeyDown(rl.KeyE) {
		rotateAngle := -rotateSpeed * rl.GetFrameTime()
		camera.Target = rotateAround(camera.Target, camera.Position, camera.Up, rotateAngle)
	}
}

// rotateAround rotates a vector around an axis by a given angle.
func rotateAround(target, position, axis rl.Vector3, angle float32) rl.Vector3 {
	rotationMatrix := rl.MatrixRotate(axis, angle)
	direction := rl.Vector3Subtract(target, position)
	direction = rl.Vector3Transform(direction, rotationMatrix)
	return rl.Vector3Add(position, direction)
}
