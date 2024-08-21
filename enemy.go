package main

import rl "github.com/gen2brain/raylib-go/raylib"

type EnemyType = int32

const (
	Slime EnemyType = iota
	NotSlime
)

type Intention = int32

const (
	Neutral Intention = iota
	Offensive
	Defansive
)

type Enemy struct {
	TexturePath      string
	AttackRange      []float32
	Health           float32
	Defence          float32
	Type             EnemyType
	CurrentIntention Intention
	Texture          rl.Texture2D
}

// TODO: load texture on setup()
func DefineEnemy(enemyType EnemyType, atkRange []float32, health, defence float32, path string) Enemy {
	return Enemy{
		Type:             enemyType,
		Health:           health,
		AttackRange:      atkRange,
		Defence:          defence,
		CurrentIntention: Neutral,
		Texture:          rl.LoadTexture(path),
	}
}

func (node *Node) assingEnemy(enemy Enemy) {
	node.Spawner.Enemy = &enemy
}

func (node *Node) killEnemy() {
	node.Spawner.Enemy = nil
}

func (node *Node) SetSpawner(enemyType EnemyType) {
	e, found := Scene.Enemies[enemyType]
	if !found {
		panic("Enemy type not defined")
	}
	node.assingEnemy(e)
}

func (e *Enemy) attackPlayer(p *Player, dmg float32) {
	var unblocked float32 = 0
	if p.Defence-dmg >= 0 {
		p.Defence -= dmg
	} else {
		unblocked = dmg - p.Defence
		p.Defence = 0
	}
	p.Health -= unblocked

	if e.Health <= 0 {
		// TODO:  handle player death
	}
}

func (e *Enemy) block() {}
