package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (s *Scene3D) render() {
	rl.ClearBackground(rl.DarkGray)
	rl.BeginMode3D(*GameState.camera)

	// rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(32, 32), rl.DarkBrown)

	maze.draw()

	for _, item := range s.Items {
		if item.hidden {
			continue
		}

		var tint color.RGBA
		if item.highlight {
			tint = rl.Purple
		} else {
			tint = rl.White
		}

		if GameState.editFocusedItemUid != "" && item.uid != GameState.editFocusedItemUid {
			tint = rl.NewColor(0, 0, 0, 70)
		}
		rl.DrawModel(item.model, item.pos, item.scale, tint)
	}

	// rl.DrawGrid(20, 8)
	rl.EndMode3D()
}
