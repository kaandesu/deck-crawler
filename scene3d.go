package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (s *Scene3D) render() {
	rl.ClearBackground(rl.SkyBlue)
	rl.BeginMode3D(*GameState.camera)

	// maze.draw()

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
		size := maze.scale * float32(len(maze.matrix))
		rl.DrawPlane(rl.NewVector3(size/2-8, -0.0001, size/2-8), rl.NewVector2(size, size), rl.Gray)
	}

	for _, node := range maze.allNodes {
		if node.Spawner.Enemy != nil {
			rl.DrawBillboard(*GameState.camera, node.Spawner.Enemy.Texture, rl.NewVector3(node.posX, 3, node.posY), 6.0, rl.White)
		}
	}

	rl.EndMode3D()
}
