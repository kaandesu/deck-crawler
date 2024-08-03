- basic ui design
- creating the borders between the different views (3d view, deck, stats etc)
- world's simplest main menu

NOTES:

- getting 3 cards from the deck at each round
- cards can have:
  - temp effects for the rest of the rounds
  - can have immidiate effects
  - can have conditional effects
- basic cards to implement:
  - simple atack
  - simple defence
  - stun an enemy
  - run from enemy
  - (future next) try to dodge attack with a percentage (depends on the level diff)
- see the intentions of the enemy, enemy can:

  - attack
  - defence
  - get stunned for a round

  CURRENT TODOS:

  - make state,style,viewport all part of the "game engine"
  - save the state of the changes from editor mode: but this implies
    saving stuff to some json and we get the state of each and every assets in the scene
    maybe like (levels/level1/wall1.json etc)
