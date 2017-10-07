package collision

import "github.com/jakecoffman/cp"

const (
	Wall   cp.CollisionType = iota << 0
	Ship
	Bullet
	Shield
)
