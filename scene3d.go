package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ground   rl.Model
	testCube rl.Model
)

func (s *Scene3D) render() {
	rl.ClearBackground(rl.DarkBlue)
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
		item.model.Materials.Shader = shader
		rl.DrawModel(item.model, item.pos, item.scale, tint)
		size := maze.scale * float32(len(maze.matrix))

		rl.DrawModel(ground, rl.NewVector3(size/2-8, -0.0001, size/2-8), size, rl.Gray)
		rl.DrawModel(testCube, rl.NewVector3(0, 4, 55), 2, rl.White)

		for i := 0; i < len(lights); i++ {
			if lights[i].enabled == 1 {
				rl.DrawSphereEx(lights[i].position, 0.2, 8, 8, lights[i].color)
			} else {
				rl.DrawSphereWires(lights[i].position, 0.2, 8, 8, rl.Fade(lights[i].color, 0.3))
			}
		}

	}
	cameraPos := []float32{GameState.camera.Position.X, GameState.camera.Position.Y, GameState.camera.Position.Z}
	rl.SetShaderValue(shader, *shader.Locs, cameraPos, rl.ShaderUniformVec3)

	for _, node := range maze.allNodes {
		if node.Spawner.Enemy != nil {
			rl.DrawBillboard(*GameState.camera, node.Spawner.Enemy.Texture, rl.NewVector3(node.posX, 3, node.posY), 6.0, rl.White)
		}
	}

	rl.EndMode3D()
}
