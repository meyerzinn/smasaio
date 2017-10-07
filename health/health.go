package health

import (
	"github.com/20zinnm/smasaio/ecs"
	"github.com/20zinnm/smasaio/component"
)

type System struct {
	entities []entity
	world    *ecs.World
}

func (s *System) Update(dt float64) {
	for _, entity := range s.entities {
		if entity.Health.Health <= 0 {
			go s.world.RemoveEntity(entity.EntityID)
		}
	}
}

func (s *System) Add(id ecs.EntityID, health *component.Health) {
	s.entities = append(s.entities, entity{id, health})
}

func (s *System) Remove(id ecs.EntityID) {
	var delete int = -1
	for index, e := range s.entities {
		if e.EntityID == id {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

type entity struct {
	ecs.EntityID
	*component.Health
}
