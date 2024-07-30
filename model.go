package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	BoxModel      *rl.Model = &rl.Model{}
	WallModel     *rl.Model = &rl.Model{}
	WallDoorModel *rl.Model = &rl.Model{}
)

func LoadModels() {
	boxModel := rl.LoadModel("./res/gltf/box_large.gltf")
	wallModel := rl.LoadModel("./res/gltf/wall.gltf")
	wallDoorModel := rl.LoadModel("./res/gltf/wall_doorway.gltf")

	WallModel = &wallModel
	WallDoorModel = &wallDoorModel
	BoxModel = &boxModel

	WallModel.Transform = rl.MatrixRotateXYZ(rl.NewVector3(0, 90, 0))
}
