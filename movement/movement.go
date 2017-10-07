package movement

import (
	"github.com/20zinnm/smasaio/component"
	"github.com/jakecoffman/cp"
	"github.com/20zinnm/smasaio/ecs"
)

type entity struct {
	ecs.EntityID
	*component.Input
	*component.Physics
}

type System struct {
	entities []entity
}

func NewSystem() ecs.System {
	return &System{}
}

func (s *System) Update(dt float64) {
	for i := 0; i < len(s.entities); i++ {
		if s.entities[i].Left {
			if !s.entities[i].Right { // both = neither {
				s.entities[i].Body.SetAngularVelocity(-0.5)
			}
		} else if s.entities[i].Right {
			s.entities[i].Body.SetAngularVelocity(0.5)
		}
		if s.entities[i].Thrusting {
			s.entities[i].Body.ApplyImpulseAtLocalPoint(cp.Vector{10, 10}, cp.Vector{0, 0})
		}
	}
}

func (s *System) Remove(id ecs.EntityID) {
	var delete int = -1
	for index, entity := range s.entities {
		if entity.EntityID == id {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}
