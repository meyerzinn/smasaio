package entities

import (
	"github.com/20zinnm/smasaio/ecs"
	"github.com/20zinnm/smasaio/collision"
	"github.com/jakecoffman/cp"
	"github.com/20zinnm/smasaio/component"
	"github.com/20zinnm/smasaio/health"
)

func NewBullet(world *ecs.World, owner ecs.EntityID, lifespan int) ecs.EntityID {
	id := world.NewEntity()

	chealth := component.Health{
		Current: 100, Max: 100,
	}

	cbullet := component.Bullet{
		Owner:     owner,
		TicksLeft: lifespan,
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *collision.System:
			var ownerBody *cp.Body
			sys.Space.EachBody(func(b *cp.Body) {
				if bid, ok := b.UserData.(ecs.EntityID); ok {
					if bid == owner {
						ownerBody = b
					}
				}
			})
			if ownerBody == nil {
				panic("bullet must have an owner")
			}
			phys := newBulletPhysics(sys.Space, id, ownerBody)
			sys.Add(id, phys, nil, nil, &chealth, &cbullet)
		case *health.System:
			sys.Add(id, &chealth)
		}
	}

	return id
}

func newBulletPhysics(space *cp.Space, id ecs.EntityID, owner *cp.Body) *component.Physics {
	body := cp.NewBody(0, 0)
	shell := cp.NewCircle(body, 3, cp.Vector{0, 0})
	shell.SetCollisionType(collision.Bullet)
	body.AddShape(shell)
	body.SetPosition(owner.Position())
	body.SetVelocityVector(owner.Velocity().Add(owner.LocalToWorld(cp.Vector{0, 30})))
	body.UserData = id
	space.AddBody(body)
	return &component.Physics{Body: body}
}
