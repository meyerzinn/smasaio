package shooting

import (
	"github.com/20zinnm/smasaio/ecs"
	"github.com/20zinnm/smasaio/component"
	"github.com/20zinnm/smasaio/entities"
)

type System struct {
	entities []entity
	world    *ecs.World
}

type entity struct {
	ecs.EntityID
	*component.Input
	*component.Cannon
}

func (s *System) Update(dt float64) {
	for index, entity := range s.entities {
		if entity.Shooting && entity.Cooldown == 0 {
			entities.NewBullet(s.world, entity.EntityID, s.entities[index].BulletLifetime)
			s.entities[index].Cooldown = 250
		} else if entity.Cooldown > 0 {
			entity.Cooldown--
		}
	}
}

func (s *System) Add(id ecs.EntityID, input *component.Input, cannon *component.Cannon) {
	s.entities = append(s.entities, entity{id, input, cannon})
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
