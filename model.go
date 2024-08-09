package main

import (
	"errors"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ModelType int32

const (
	Wall ModelType = iota
	WallDoorway
	Box
)

func (scene *Scene3D) AddModel(modelType ModelType, pos, rot rl.Vector3, scale float32) error {
	var (
		found = false
		path  = ""
	)

	uid := fmt.Sprintf("wall%d", len(scene.Items))
	if _, e := scene.Items[uid]; e {
		return errors.New("model.go: same key already exists: " + uid)
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
	scene.Items[uid] = &SceneItem{
		uid:   uid,
		model: temp,
		pos:   pos,
		rot:   rot,
		scale: scale,
	}
	return nil
}
