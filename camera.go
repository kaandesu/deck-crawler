package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera = rl.Camera3D

func NewCamera() *rl.Camera {
	// Define the initial position
	initialPosition := rl.NewVector3(15.0, 5.0, 15.0)
	target := rl.NewVector3(0.0, 0.0, 0.0)
	up := rl.NewVector3(0.0, 1.0, 0.0)

	// Rotate the initial position by 45 degrees around the Y-axis (up vector)
	var rotationAngle float32 = 45.0 * rl.Deg2rad
	rotatedPosition := rotateAround(initialPosition, target, up, rotationAngle)

	return &rl.Camera{
		Position:   rotatedPosition,
		Target:     target,
		Up:         up,
		Fovy:       50.0,
		Projection: rl.CameraPerspective,
	}
}

func MoveCam(cam *Camera, direction rl.Vector3) {
	cam.Position = rl.Vector3Add(cam.Position, direction)
	// cam.Target = rl.Vector3Add(cam.Target, direction)
}

func UpdateCameraCustom(camera *rl.Camera) {
	var moveSpeed float32 = 0.2

	forwardDir := rl.Vector3Subtract(camera.Target, camera.Position)
	forwardDir.Y = 0
	forwardDir = rl.Vector3Normalize(forwardDir)

	if rl.IsKeyDown(rl.KeyW) {
		camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(forwardDir, moveSpeed))
		camera.Target = rl.Vector3Add(camera.Target, rl.Vector3Scale(forwardDir, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		camera.Position = rl.Vector3Subtract(camera.Position, rl.Vector3Scale(forwardDir, moveSpeed))
		camera.Target = rl.Vector3Subtract(camera.Target, rl.Vector3Scale(forwardDir, moveSpeed))
	}

	// Handle rotation
	if turningLeft || turningRight {
		rotateStep := rotationSpeed * rl.GetFrameTime()
		if turningRight {
			rotateStep = -rotateStep
		}

		currentRotation += rotateStep
		if float32(math.Abs(float64(currentRotation))) >= float32(math.Abs(float64(targetRotation))) {
			rotateStep -= currentRotation - targetRotation
			turningLeft, turningRight = false, false
			inputBlocked = false
		}

		camera.Target = rotateAround(camera.Target, camera.Position, camera.Up, rotateStep*rl.Deg2rad)
	}
}

// rotateAround rotates a vector around an axis by a given angle.
func rotateAround(target, position, axis rl.Vector3, angle float32) rl.Vector3 {
	rotationMatrix := rl.MatrixRotate(axis, angle)
	direction := rl.Vector3Subtract(target, position)
	direction = rl.Vector3Transform(direction, rotationMatrix)
	return rl.Vector3Add(position, direction)
}
