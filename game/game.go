package game

import (
	"time"
	"github.com/20zinnm/smasaio/collision"
	"github.com/jakecoffman/cp"
	"github.com/20zinnm/smasaio/movement"
	"github.com/20zinnm/smasaio/ecs"
)

type Game struct {
	world    *ecs.World
	stop     chan struct{}
	tickRate time.Duration
}

func (g *Game) tick(dt float64) {
	g.world.Update(dt)
}

func (g *Game) Stop() {
	g.stop <- struct{}{}
}

type Option func(game *Game)

func NewGame(options ...Option) {
	game := Game{
		world:    &ecs.World{},
		stop:     make(chan struct{}),
		tickRate: 50 * time.Millisecond,
	}
	for _, o := range options {
		o(&game)
	}

	space := newSpace()
	space.SetGravity(cp.Vector{0, 0})
	game.world.AddSystem(movement.NewSystem())
	game.world.AddSystem(collision.NewSystem())

	go func() {
		tickInterval := time.Second / game.tickRate
		timeStart := time.Now().UnixNano()

		ticker := time.NewTicker(tickInterval)

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()
				delta := float64(now-timeStart) / 1000000000
				timeStart = now
				game.tick(delta)
			case <-game.stop:
				ticker.Stop()
			}
		}
	}()
}

func newSpace() *cp.Space {
	space := cp.NewSpace()

	return space
}
