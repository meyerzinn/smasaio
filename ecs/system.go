package ecs

type System interface {
	// Update updates the system. It is invoked by the engine once every frame,
	// with dt being the duration since the previous update.
	Update(dt float64)

	// Remove removes the given entity from the system.
	Remove(id EntityID)
}