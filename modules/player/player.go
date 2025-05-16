package player

import (
	"fmt"
	"meermookh/modules/aabb"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	rect       rl.Rectangle
	speed      uint
	jumpHeight float32

	canJump    bool
	isStanding bool
	isJumping  bool

	mu sync.RWMutex
}

func New(pos rl.Vector2) Player {
	return Player{
		jumpHeight: 150,
		speed:      5,
		isStanding: false,
		isJumping:  false,
		rect: rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  32,
			Height: 32,
		},
	}
}

func (p *Player) GetRect() rl.Rectangle {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.rect
}

func (p *Player) Draw() {
	origin := rl.Vector2{
		X: float32(p.rect.Width) / 2,
		Y: float32(p.rect.Height) / 2,
	}

	color := rl.Color{
		R: 128,
		G: 128,
		B: 128,
		A: 255,
	}

	rl.DrawRectanglePro(p.rect, origin, 0, color)
}

func (p *Player) Update() {
	p.mu.RLock()
	canJumpNow := p.canJump && !p.isJumping
	p.mu.RUnlock()

	p.mu.Lock()
	if !p.isStanding {
		p.rect.Y += float32(p.speed)
	}

	if rl.IsKeyDown(rl.KeyA) {
		p.rect.X -= float32(p.speed)
	}

	if rl.IsKeyDown(rl.KeyD) {
		p.rect.X += float32(p.speed)
	}
	p.mu.Unlock()

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) || rl.IsKeyPressed(rl.KeyLeftControl) {
		go p.attack()
	}

	if rl.IsKeyDown(rl.KeySpace) && canJumpNow {
		go p.jump()
	}
}

func (p *Player) HandleCollision(info aabb.CollisionInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if info.IsCollided {
		p.isStanding = info.IsStanding
		if info.IsStanding {
			p.canJump = true
			p.isJumping = false
		}
	}
}

func (p *Player) ResetCollision() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.isStanding = false
	p.canJump = false
}

func (p *Player) jump() {
	p.mu.Lock()
	p.isJumping = true
	p.canJump = false
	p.mu.Unlock()

	var jumpedAt float32 = 0.0
	var jumpAmpl float32 = 0.005

	const updateInterval = 5
	var accumulatedChange float32 = 0

	for jumpedAt <= p.jumpHeight {
		accumulatedChange -= float32(jumpAmpl)
		jumpedAt += jumpAmpl

		if int(jumpedAt/float32(jumpAmpl))%updateInterval == 0 {
			p.mu.Lock()
			p.rect.Y += accumulatedChange
			p.mu.Unlock()
			accumulatedChange = 0
		}

		fmt.Printf("jumpedAt: %f\n", jumpedAt)
	}
}

func (p *Player) attack() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// todo: attack
}
