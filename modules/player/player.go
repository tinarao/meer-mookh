package player

import (
	"meermookh/config"
	"meermookh/modules/aabb"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	rect       rl.Rectangle
	speed      uint
	isStanding bool
}

func New() Player {
	return Player{
		speed:      5,
		isStanding: false,
		rect: rl.Rectangle{
			X:      100,
			Y:      100,
			Width:  32,
			Height: 32,
		},
	}
}

func (p *Player) GetRect() rl.Rectangle {
	return p.rect
}

func (p *Player) Draw() {
	origin := rl.Vector2{
		X: float32(config.BASE_TILE_SIZE) / 2,
		Y: float32(config.BASE_TILE_SIZE) / 2,
	}

	rl.DrawRectanglePro(p.rect, origin, 0, rl.Red)
}

func (p *Player) Update() {
	if !p.isStanding {
		p.rect.Y += float32(p.speed)
	}

	if rl.IsKeyDown(rl.KeyA) {
		p.rect.X -= float32(p.speed)
	}

	if rl.IsKeyDown(rl.KeyD) {
		p.rect.X += float32(p.speed)
	}
}

func (p *Player) HandleCollision(info aabb.CollisionInfo) {
	if info.IsCollided {
		p.isStanding = info.IsStanding
	}
}
