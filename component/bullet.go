package component

import "github.com/20zinnm/smasaio/ecs"

type Bullet struct {
	Owner     ecs.EntityID
	TicksLeft int
}
