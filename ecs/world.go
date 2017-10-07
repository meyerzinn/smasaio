package ecs

import (
	"sync/atomic"
	"sync"
)

type World struct {
	systems []System
	lock    *sync.RWMutex
	max     EntityID
}

// AddSystem adds the given System to the World, sorted by priority.
func (w *World) AddSystem(system System) {
	w.lock.Lock()
	w.systems = append(w.systems, system)
	w.lock.Unlock()
}

// Systems returns the list of Systems managed by the World.
func (w *World) Systems() []System {
	return w.systems
}

// Update updates each System managed by the World. It is invoked by the engine
// once every frame, with dt being the duration since the previous update.
func (w *World) Update(dt float64) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	for _, system := range w.Systems() {
		system.Update(dt)
	}
}

// RemoveEntity removes the entity across all systems.
func (w *World) RemoveEntity(id EntityID) {
	w.lock.Lock()
	for _, sys := range w.systems {
		sys.Remove(id)
	}
	w.lock.Unlock()
}

func (w *World) NewEntity() EntityID {
	return EntityID(atomic.AddUint64(&w.max, 1) - 1)
}
