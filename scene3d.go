package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (s *Scene3D) render() {
	rl.ClearBackground(rl.SkyBlue)
	rl.BeginMode3D(*GameState.camera)

	rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(32, 32), rl.DarkBrown)

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
		// FIXME: add focus
		rl.DrawModel(item.model, item.pos, item.scale, tint)
	}

	rl.DrawGrid(50, 0)
	rl.EndMode3D()
}