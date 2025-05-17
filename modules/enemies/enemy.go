package enemies

import (
	"meermookh/config"
	"meermookh/modules/aabb"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Enemy struct {
	rect       rl.Rectangle
	hp         float32
	isStanding bool
	speed      float32
	gravity    float32
	color      rl.Color

	mu sync.Mutex
}

func New(pos rl.Vector2) Enemy {
	return Enemy{
		isStanding: false,
		speed:      4,
		hp:         100,
		gravity:    5,
		color:      rl.Red,
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

	rl.DrawRectanglePro(e.rect, origin, 0, e.color)
}

func (e *Enemy) Update(tiles *[]aabb.Drawable) {
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

func (e *Enemy) SetColor(color rl.Color) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.color = color
}

func (e *Enemy) ApplyDamage(damage float32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.hp -= damage
}

func (e *Enemy) GetHP() float32 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.hp
}
