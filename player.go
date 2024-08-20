package main

type CardType = int32

const (
	Attack CardType = iota
	Block
)

// TODO: add the deck
type Player struct {
	Attack  float32
	Health  float32
	Defence float32
	Energy  int32
}

func (node *Node) attackEnemy(e *Enemy, dmg float32) {
	var unblocked float32 = 0
	if e.Defence-dmg >= 0 {
		e.Defence -= dmg
	} else {
		unblocked = dmg - e.Defence
		e.Defence = 0
	}
	e.Health -= unblocked
	if e.Health <= 0 {
		node.killEnemy()
	}
}
