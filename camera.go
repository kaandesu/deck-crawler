package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Camera = rl.Camera3D

func NewCamera() *Camera {
	return &Camera{
		Position: rl.NewVector3(10.0, 10.0, 10.0),
		Target:   rl.NewVector3(0.0, 0.0, 0.0),
		Up:       rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:     45.0,
	}
}

func MoveCam(cam *Camera, direction rl.Vector3) {
	cam.Position = rl.Vector3Add(cam.Position, direction)
	cam.Target = rl.Vector3Add(cam.Target, direction)
}

func UpdateCam(cam *Camera, deltaTime float32) {
	moveSpeed := deltaTime * 1.0
	front := rl.Vector3Subtract(cam.Target, cam.Position)
	front = rl.Vector3Normalize(front)
	front.X *= moveSpeed
	front.Z *= moveSpeed

	right := rl.Vector3CrossProduct(cam.Up, front)
	right = rl.Vector3Normalize(right)
	right.X *= moveSpeed
	right.Z *= moveSpeed

	if rl.IsKeyDown(rl.KeyW) {
		MoveCam(cam, front)
	}
	if rl.IsKeyDown(rl.KeyS) {
		MoveCam(cam, rl.Vector3Negate(front))
	}
	if rl.IsKeyDown(rl.KeyA) {
		MoveCam(cam, right)
	}
	if rl.IsKeyDown(rl.KeyD) {
		MoveCam(cam, rl.Vector3Negate(right))
	}
}
