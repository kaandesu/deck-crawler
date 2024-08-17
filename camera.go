package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera = rl.Camera3D

func NewCamera() *rl.Camera {
	initialPosition := rl.NewVector3(0.01, 5.0, 0.01)
	target := rl.NewVector3(0.01, 5.0, 0.0)
	up := rl.NewVector3(0.0, 1.0, 0.0)

	var rotationAngle float32 = -180.0 * rl.Deg2rad
	rotatedPosition := rotateAround(initialPosition, target, up, rotationAngle)

	return &rl.Camera{
		Position:   rotatedPosition,
		Target:     target,
		Up:         up,
		Fovy:       60.0,
		Projection: rl.CameraPerspective,
	}
}

var targetRotation float32

var (
	targetPos    rl.Vector3
	elapsedTime  float32 = 0
	movingToNode bool    = false
)

func UpdateCameraCustom(camera *rl.Camera) {
	var moveSpeed float32 = 0.04

	forwardDir := rl.Vector3Subtract(camera.Target, camera.Position)
	forwardDir.Y = 0
	forwardDir = rl.Vector3Normalize(forwardDir)

	if movingToNode {
		movementDir := rl.Vector3Subtract(targetPos, camera.Position)
		movementDir.Y = 0
		distance := rl.Vector3Length(movementDir)

		movementDir = rl.Vector3Normalize(movementDir)
		moveStep := moveSpeed * rl.GetFrameTime() * 1000

		if distance > moveStep {
			camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(movementDir, moveStep))
			camera.Target = rl.Vector3Add(camera.Position, forwardDir)
		} else {
			camera.Position = targetPos
			camera.Target = rl.Vector3Add(camera.Position, forwardDir)
			movingForward = false
			movingToNode = false
			elapsedTime = 0
		}
	}

	if movingBackward {
		// TODO: implement here
	}

	if turningLeft || turningRight {
		rotateStep := 180 * rl.GetFrameTime()
		if turningRight {
			targetRotation = -90
			rotateStep = -rotateStep
		} else {
			targetRotation = 90
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

func rotateAround(target, position, axis rl.Vector3, angle float32) rl.Vector3 {
	rotationMatrix := rl.MatrixRotate(axis, angle)
	direction := rl.Vector3Subtract(target, position)
	direction = rl.Vector3Transform(direction, rotationMatrix)
	return rl.Vector3Add(position, direction)
}
