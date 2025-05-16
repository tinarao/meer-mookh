package enemies

import (
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/tile"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	rect       rl.Rectangle
	isStanding bool
	speed      float32
	gravity    float32

	mu sync.Mutex
}

func New(pos rl.Vector2) Enemy {
	return Enemy{
		isStanding: false,
		speed:      4,
		gravity:    5,
		rect: rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  config.BASE_TILE_SIZE,
			Height: config.BASE_TILE_SIZE,
		},
	}
}

func (e *Enemy) Draw() {
	e.mu.Lock()
	defer e.mu.Unlock()

	origin := rl.Vector2{
		X: float32(e.rect.Width) / 2,
		Y: float32(e.rect.Height) / 2,
	}

	rl.DrawRectanglePro(e.rect, origin, 0, rl.Red)
}

func (e *Enemy) Update(tiles *[]tile.Tile) {
	e.mu.Lock()
	defer e.mu.Unlock()

	coll := aabb.Check(&e.rect, tiles)
	if !coll.IsStanding {
		e.rect.Y += e.gravity
	}
}

func (e *Enemy) GetRect() rl.Rectangle {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.rect
}
