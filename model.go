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
		uid:   modelName,
		model: temp,
		pos:   pos,
		rot:   rot,
		scale: scale,
	}
	return nil
}

func SetupModels() {
	var (
		rot           = rl.NewVector3(0, 0, 0)
		scale float32 = 2.2
	)
	_ = rot
	_ = scale
	// TODO: add a id system, something like wall#123
	// WARN: disabled models
	/*
		Scene.AddModel(WallDoorway, "door1", rl.NewVector3(0, 0, -10), rot, scale)
		Scene.AddModel(Wall, "wall1", rl.NewVector3(-3.3, 0, -5.5), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
		Scene.AddModel(Wall, "wall2", rl.NewVector3(3.4, 0, -5.5), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
		Scene.AddModel(Wall, "wall3", rl.NewVector3(-3.3, 0, 3.3), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
		Scene.AddModel(Wall, "wall4", rl.NewVector3(3.4, 0, 3.3), rl.NewVector3(0, 90*rl.Deg2rad, 0), scale)
	*/
}

func LoadModels() {
	SetupModels()
}
