package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	BoxModel                = rl.Model{}
	WallModel     *rl.Model = &rl.Model{}
	WallDoorModel           = rl.Model{}
)

func LoadModels() {
	BoxModel = rl.LoadModel("./res/gltf/box_large.gltf")
	temp := rl.LoadModel("./res/gltf/wall.gltf")
	WallModel = &temp
	WallDoorModel = rl.LoadModel("./res/gltf/wall_doorway.gltf")

	WallModel.Transform = rl.MatrixRotateXYZ(rl.NewVector3(0, 90, 0))
	// WallModel.Transform = rl.MatrixRotateXYZ(rl.NewVector3(90, 45, 10))
}
