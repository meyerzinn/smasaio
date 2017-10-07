package entity

import (
	"github.com/gorilla/websocket"
	"github.com/20zinnm/smasaio/ecs"
)

func NewPlayer(world *ecs.World, conn *websocket.Conn) ecs.EntityID {
	entity := ecs.NewBasic()

	return entity
	//for _, system := range world.Systems() {
	//	switch system.(type) {
	//	case networking.System:
	//		system.(*networking.System).Add(entity, conn)
	//	}
	//}
}
