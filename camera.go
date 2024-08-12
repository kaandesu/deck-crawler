package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera = rl.Camera3D

func NewCamera() *rl.Camera {
	node := maze.matrix[0][0]
	fmt.Printf(">>>%+v \n", node)
	// initialPosition := rl.NewVector3(15+node.posX, 45.0, 15+node.posY)
	initialPosition := rl.NewVector3(15, 5.0, 15)
	target := rl.NewVector3(0.0, 0.0, 0.0)
	up := rl.NewVector3(0.0, 1.0, 0.0)

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

var targetRotation float32

func UpdateCameraCustom(camera *rl.Camera) {
	var moveSpeed float32 = 0.2

	forwardDir := rl.Vector3Subtract(camera.Target, camera.Position)
	forwardDir.Y = 0
	forwardDir = rl.Vector3Normalize(forwardDir)

	if movingForward {
		camera.Position = rl.Vector3Add(camera.Position, rl.Vector3Scale(forwardDir, moveSpeed))
		camera.Target = rl.Vector3Add(camera.Target, rl.Vector3Scale(forwardDir, moveSpeed))
	}
	if movingBackward {
		camera.Position = rl.Vector3Subtract(camera.Position, rl.Vector3Scale(forwardDir, moveSpeed))
		camera.Target = rl.Vector3Subtract(camera.Target, rl.Vector3Scale(forwardDir, moveSpeed))
	}

	if turningLeft || turningRight {
		rotateStep := 90 * rl.GetFrameTime()
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
