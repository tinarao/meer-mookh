package player

import (
	"fmt"
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/tile"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	rect       rl.Rectangle
	speed      uint
	hp         int
	jumpHeight float32

	canJump    bool
	isStanding bool
	isJumping  bool

	mu sync.RWMutex

	vel                  rl.Vector2
	isFallingBelowScreen bool
	fallDamageTicker     *time.Ticker
	fallDamageDone       chan struct{}
}

func New(pos rl.Vector2) Player {
	return Player{
		jumpHeight: 150,
		speed:      5,
		hp:         100,
		isStanding: false,
		isJumping:  false,
		rect: rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  config.PLAYER_WIDTH,
			Height: config.PLAYER_HEIGHT,
		},
		vel:            rl.Vector2{X: 0, Y: 0},
		fallDamageDone: make(chan struct{}),
	}
}

func (p *Player) GetRect() rl.Rectangle {
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

func (p *Player) Update(tiles *[]tile.Tile) {
	p.mu.Lock()
	canJumpNow := p.canJump && !p.isJumping
	p.mu.Unlock()

	p.HandleCollision(tiles)

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

	p.mu.Lock()
	if p.rect.Y >= config.WINDOW_H {
		if !p.isFallingBelowScreen {
			p.isFallingBelowScreen = true
			p.fallDamageTicker = time.NewTicker(time.Millisecond * 75)
			go func() {
				for {
					select {
					case <-p.fallDamageTicker.C:
						p.DealDamage(5)
						fmt.Printf("Damage. Current hp: %d\n", p.GetHP())
					case <-p.fallDamageDone:
						return
					}
				}
			}()
		}
	} else {
		if p.isFallingBelowScreen {
			p.isFallingBelowScreen = false
			if p.fallDamageTicker != nil {
				p.fallDamageTicker.Stop()
				close(p.fallDamageDone)
				p.fallDamageDone = make(chan struct{})
			}
		}
	}
	p.mu.Unlock()
}

func (p *Player) HandleCollision(tiles *[]tile.Tile) {
	plRect := p.GetRect()
	info := aabb.Check(&plRect, tiles)
	if info.IsCollided {
		p.mu.Lock()
		p.isStanding = info.IsStanding
		if info.IsStanding {
			p.canJump = true
			p.isJumping = false
		}
		p.mu.Unlock()
	} else {
		p.ResetCollision()
	}
}

func (p *Player) ResetCollision() {
	p.mu.Lock()
	p.isStanding = false
	p.canJump = false
	p.mu.Unlock()
}

func (p *Player) GetHP() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.hp
}

func (p *Player) DealDamage(amount int) {
	p.mu.Lock()
	p.hp -= amount
	if p.hp < 0 {
		p.hp = 0
	}
	p.mu.Unlock()
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
