package collision

import (
	"github.com/jakecoffman/cp"
	"github.com/20zinnm/smasaio/component"
	"github.com/20zinnm/smasaio/ecs"
)

type System struct {
	entities []entity
	Space    *cp.Space
	world    *ecs.World
}

type entity struct {
	ecs.EntityID
	Physics *component.Physics
	//Input   *component.Input
	Health *component.Health
	Shield *component.Shield // only present on ships
	Bullet *component.Bullet // only present on bullets
}

func (s *System) Add(id ecs.EntityID, physics *component.Physics, health *component.Health, /*input *component.Input,*/ shield *component.Shield, bullet *component.Bullet) {
	s.entities = append(s.entities, entity{EntityID: id, Physics: physics, /*Input: input,*/ Health: health, Shield: shield, Bullet: bullet})
}

func (s *System) Update(dt float64) {
	// Advance simulation
	s.Space.Step(dt)

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

func NewSystem() ecs.System {
	//system := System{
	//lock: &sync.RWMutex{},
	//}
	var system System

	// initialize space
	space := cp.NewSpace()
	space.NewCollisionHandler(Wall, Bullet).BeginFunc = collideWallBullet(&system)
	space.NewCollisionHandler(Bullet, Ship).BeginFunc = collideBulletShip(&system)
	space.NewCollisionHandler(Bullet, Bullet).BeginFunc = collideBulletBullet(&system)
	space.NewCollisionHandler(Bullet, Shield).BeginFunc = collideBulletShield(&system)

	/* Add boundaries */
	// left
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{-1000, -1000}, cp.Vector{-1000, 1000}, 0))
	// top
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{-1000, 1000}, cp.Vector{1000, 1000}, 0))
	// right
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{1000, 1000}, cp.Vector{1000, -1000}, 0))
	// bottom
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{1000, -1000}, cp.Vector{-1000, -1000}, 0))
	system.Space = space

	return &system
}

func collideWallBullet(s *System) cp.CollisionBeginFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		_, b := arb.Bodies()
		id := b.UserData.(ecs.EntityID)
		for index, entity := range s.entities {
			if entity.EntityID == id {
				if s.entities[index].Health != nil {
					s.entities[index].Health.Current -= 20
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
		shipid := playerb.UserData.(ecs.EntityID)
		var bi, si int
		for index, entity := range s.entities {
			switch entity.EntityID {
			case bulletid:
				bi = index
			case shipid:
				si = index
			}
		}
		if s.entities[si].EntityID == s.entities[bi].Bullet.Owner {
			return arb.Ignore()
		}
		s.entities[si].Health.Current -= s.entities[bi].Health.Current
		s.entities[bi].Health.Current = 0
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
		dmg1, dmg2 := s.entities[b2i].Health.Current, s.entities[b1i].Health.Current
		s.entities[b1i].Health.Current -= dmg1
		s.entities[b2i].Health.Current -= dmg2
		return true
	}
}

func collideBulletShield(s *System) cp.CollisionBeginFunc {
	return func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		bulletb, shieldb := arb.Bodies()
		bulletid := bulletb.UserData.(ecs.EntityID)
		shipid := shieldb.UserData.(ecs.EntityID)
		var bi, si int
		for index, entity := range s.entities {
			switch entity.EntityID {
			case bulletid:
				bi = index
			case shipid:
				si = index
			}
		}
		if s.entities[si].Shield.Active {
			s.entities[bi].Health.Current = 0
		}
		return true
	}
}
