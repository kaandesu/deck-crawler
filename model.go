package main

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ModelType int32

const (
	Wall ModelType = iota
	WallDoorway
	Box
)

func (scene *Scene3D) AddModel(modelType ModelType, modelName string, pos, rot rl.Vector3, scale float32) error {
	var (
		found = false
		path  = ""
	)

	if _, e := scene.Items[modelName]; e {
		return errors.New("model.go: same key already exists: " + modelName)
	}
	switch modelType {
	case Wall:
		path = "./res/gltf/wall.gltf"
		found = true
	case WallDoorway:
		path = "./res/gltf/wall_doorway.gltf"
		found = true
	case Box:
		path = "./res/gltf/box_large.gltf"
		found = true
	}
	if !found {
		return errors.New("model key not found")
	}

	temp := rl.LoadModel(path)
	temp.Transform = rl.MatrixRotateXYZ(rot)
	scene.Items[modelName] = &SceneItem{
		model: temp,
		pos:   pos,
		rot:   rot,
		scale: scale,
	}
	return nil
}

func SetupModels() {
	var (
		z         = rl.NewVector3(0, 0, 0)
		s float32 = 2.2
	)
	ViewportState.AddModel(WallDoorway, "door1", rl.NewVector3(0, -3, -10), z, s)
	// TODO: add a id system, something like wall#123
	ViewportState.AddModel(Wall, "wall1", rl.NewVector3(-4, 0, 0), z, s)
	ViewportState.AddModel(Wall, "wall2", rl.NewVector3(4, -3, -7), rl.NewVector3(0, 1.5, 0), s)
}

func LoadModels() {
	SetupModels()
}
