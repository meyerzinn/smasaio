package collision

import (
	"github.com/jakecoffman/cp"
	"github.com/20zinnm/smasaio/component"
	"github.com/20zinnm/smasaio/ecs"
)

type System struct {
	entities []entity
	//lock     *sync.RWMutex
	space *cp.Space
	world *ecs.World
}

type entity struct {
	ecs.EntityID
	Physics component.Physics
	Input   *component.Input
	Health  *component.Health
}

//func (s *System) Add(basic *ecs.BasicEntity, ) {
//
//}

func (s *System) Update(dt float64) {
	//s.lock.RLock()
	//defer s.lock.RUnlock()
	// Advance simulation
	s.space.Step(dt)

}

func (s *System) Remove(id ecs.EntityID) {
	//s.lock.Lock()
	//defer s.lock.Unlock()

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

func NewSystem() ecs.System {
	//system := System{
	//lock: &sync.RWMutex{},
	//}
	var system System

	// initialize space
	space := cp.NewSpace()
	space.NewCollisionHandler(Wall, Bullet).BeginFunc = collideWallBullet(&system)
	space.NewCollisionHandler(Bullet, Ship).BeginFunc = collideBulletShip(&system)

	/* Add boundaries */
	// left
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{-1000, -1000}, cp.Vector{-1000, 1000}, 0))
	// top
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{-1000, 1000}, cp.Vector{1000, 1000}, 0))
	// right
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{1000, 1000}, cp.Vector{1000, -1000}, 0))
	// bottom
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{1000, -1000}, cp.Vector{-1000, -1000}, 0))
	system.space = space

	return &system
}

func collideWallBullet(s *System) cp.CollisionBeginFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		_, b := arb.Bodies()
		id := b.UserData.(ecs.EntityID)
		for index, entity := range s.entities {
			if entity.EntityID == id {
				if s.entities[index].Health != nil {
					s.entities[index].Health.Health -= 20
				}
			}
		}
		return true
	}
}

func collideBulletShip(s *System) cp.CollisionBeginFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		bulletb, playerb := arb.Bodies()
		bulletid := bulletb.UserData.(ecs.EntityID)
		playerid := playerb.UserData.(ecs.EntityID)
		var bi, pi int
		for index, entity := range s.entities {
			switch entity.EntityID {
			case bulletid:
				bi = index
			case playerid:
				pi = index
			}
		}
		s.entities[pi].Health.Health -= s.entities[bi].Health.Health
		s.entities[bi].Health.Health = 0
		return true
	}
}

func collideBulletBullet(s *System) cp.CollisionBeginFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		b1b, b2b := arb.Bodies()
		b1id := b1b.UserData.(ecs.EntityID)
		b2id := b2b.UserData.(ecs.EntityID)
		var b1i, b2i int
		for index, entity := range s.entities {
			switch entity.EntityID {
			case b1id:
				b1i = index
			case b2id:
				b2i = index
			}
		}
		// the damage to bullet 1 and bullet 2 respectively
		dmg1, dmg2 := s.entities[b2i].Health.Health, s.entities[b1i].Health.Health
		s.entities[b1i].Health.Health -= dmg1
		s.entities[b2i].Health.Health -= dmg2
		return true
	}
}
