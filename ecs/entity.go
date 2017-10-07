package ecs

type EntityID = uint64

func (e EntityID) ID() EntityID {
	return e
}

type Entity interface {
	ID() EntityID
}
