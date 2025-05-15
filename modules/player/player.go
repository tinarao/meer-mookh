package player

import (
	"fmt"
	"meermookh/config"
	"meermookh/modules/aabb"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	rect       rl.Rectangle
	speed      uint
	jumpHeight float32

	canJump    bool
	isStanding bool
	isJumping  bool
}

func New() Player {
	return Player{
		jumpHeight: 150,
		speed:      5,
		isStanding: false,
		isJumping:  false,
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

	if rl.IsKeyDown(rl.KeySpace) && !p.isJumping && p.canJump {
		go p.jump()
	}
}

func (p *Player) HandleCollision(info aabb.CollisionInfo) {
	if info.IsCollided {
		p.isStanding = info.IsStanding
		if info.IsStanding {
			p.canJump = true
			p.isJumping = false
		}
	}
}

func (p *Player) ResetCollision() {
	p.isStanding = false
	p.canJump = false
}

func (p *Player) jump() {
	p.isJumping = true
	p.canJump = false

	var jumpedAt float32 = 0.0
	var jumpAmpl float32 = 0.005
	for jumpedAt <= p.jumpHeight {
		p.rect.Y -= float32(jumpAmpl)
		jumpedAt += jumpAmpl

		// idk, when i delete this print,
		// it jumps not like it should
		fmt.Printf("jumpedAt: %f\n", jumpedAt)
	}
}
