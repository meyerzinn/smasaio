package entities

import (
	"github.com/20zinnm/smasaio/collision"
	"github.com/gorilla/websocket"
	"github.com/20zinnm/smasaio/ecs"
	"github.com/20zinnm/smasaio/component"
	"github.com/jakecoffman/cp"
	"github.com/20zinnm/smasaio/movement"
	"math/rand"
	"github.com/20zinnm/smasaio/shooting"
)

func NewPlayer(world *ecs.World, conn *websocket.Conn) ecs.EntityID {
	id := world.NewEntity()

	input := component.Input{}
	var phys *component.Physics
	health := component.Health{Current: 200, Max: 200}
	shield := component.Shield{}
	cannon := component.Cannon{}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *movement.System:
			if phys == nil {
				phys = &component.Physics{Body: sys.Space.AddBody(newPlayerPhysics(id, cp.Vector{rand.Float64(), rand.Float64()}).Body)}
			}
			sys.Add(id, &input, phys)
		case *collision.System:
			if phys == nil {
				phys = &component.Physics{Body: sys.Space.AddBody(newPlayerPhysics(id, cp.Vector{rand.Float64(), rand.Float64()}).Body)}
			}
			sys.Add(id, phys, &health, &shield, nil)
		case *shooting.System:
			sys.Add(id, &input, )
		}
	}

	return id
}

func newPlayerPhysics(id ecs.EntityID, position cp.Vector) *component.Physics {
	body := cp.NewBody(0, 0)
	body.AddShape(cp.NewPolyShape(body, 5, []cp.Vector{cp.Vector{0, 10}, cp.Vector{}}, cp.NewTransformIdentity(), 1))
	shield := cp.NewCircle(body, 20, cp.Vector{0, 0})
	shield.SetSensor(true)
	shield.SetCollisionType(collision.Shield) // collides with bullets only
	body.AddShape(shield)
	return &component.Physics{body}
}
