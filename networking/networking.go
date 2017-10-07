package networking

import (
	"engo.io/ecs"
	"github.com/gorilla/websocket"
	"github.com/20zinnm/smasaio/component"
	"github.com/bitly/go-notify"
	"github.com/20zinnm/smasaio/networking/packets"
	"github.com/google/flatbuffers/go"
	"sync"
	"github.com/jakecoffman/cp"
)

const (
	PacketWorldUpdate uint8 = iota
)

var builderPool = sync.Pool{
	New: func() interface{} {
		return flatbuffers.NewBuilder(1)
	},
}

type entity struct {
	*ecs.BasicEntity
	*websocket.Conn
	*component.Physics
	*component.Health
	*component.Ship
}

type System struct {
	entities []entity
}

func (s *System) Add(e *ecs.BasicEntity, conn *websocket.Conn, space *cp.Space, physics *component.Physics, health *component.Health, ship *component.Ship) {
	en := entity{e, conn, physics, health, ship,}
	s.entities = append(s.entities, en)
	go func(conn *websocket.Conn) {
		var ticker chan interface{}
		notify.Start("tick", ticker)
		for {
			<-ticker
			builder := builderPool.Get().(*flatbuffers.Builder)
			packets.WorldSnapshotStart(builder)
			space.PointQueryNearest()
			packets.WorldSnapshotStartDestroyVector(builder)

		}
	}(conn)
}

func (s *System) Update(dt float32) {

}

func (s *System) Remove(e ecs.BasicEntity) {
	panic("implement me")
}
