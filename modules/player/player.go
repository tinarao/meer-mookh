package player

import (
	"fmt"
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/enemies"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	rect         rl.Rectangle
	speed        uint
	damage       float32
	hp           int
	jumpHeight   float32
	frags        uint
	attackRadius int

	canJump    bool
	isStanding bool
	isJumping  bool

	mu sync.RWMutex

	vel                  rl.Vector2
	isFallingBelowScreen bool
	fallDamageTicker     *time.Ticker
	fallDamageDone       chan struct{}

	isAttacking       bool
	attackCircleTimer *time.Timer

	// kolhozny DI
	enemies *[]*enemies.Enemy
}

func New(pos rl.Vector2, enemies *[]*enemies.Enemy) Player {
	return Player{
		jumpHeight:   150,
		speed:        5,
		attackRadius: 64,
		enemies:      enemies,
		damage:       35,
		hp:           100,
		isStanding:   false,
		isJumping:    false,
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
	hpStr := fmt.Sprintf("HP: %d\n", p.hp)
	color := rl.Black
	if p.hp <= 25 {
		color = rl.Red
	}

	fragsStr := fmt.Sprintf("Frags: %d\n", p.GetFrags())
	rl.DrawText(hpStr, 50, 50, 20, color)
	rl.DrawText(fragsStr, 50, 75, 20, color)

	origin := rl.Vector2{
		X: float32(p.rect.Width) / 2,
		Y: float32(p.rect.Height) / 2,
	}

	playerColor := rl.Color{
		R: 128,
		G: 128,
		B: 128,
		A: 255,
	}

	rl.DrawRectanglePro(p.rect, origin, 0, playerColor)
}

func (p *Player) Update(tiles *[]aabb.Drawable) {
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

func (p *Player) HandleCollision(tiles *[]aabb.Drawable) {
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

func (p *Player) GetFrags() uint {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.frags
}

func (p *Player) DealDamage(amount int) {
	p.mu.Lock()
	p.hp -= amount
	if p.hp < 0 {
		p.hp = 0
	}
	p.mu.Unlock()
}

func (p *Player) GetDamage() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.damage
}

func (p *Player) AddFrags(amount uint) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.frags += amount
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
	if p.isAttacking {
		p.mu.Unlock()
		return
	}

	p.attackCircleTimer = time.NewTimer(300 * time.Millisecond)
	p.isAttacking = true
	p.mu.Unlock()

	v := rl.Vector2{
		X: p.GetRect().X,
		Y: p.GetRect().Y,
	}

	if p.enemies == nil {
		return
	}

	for _, enemy := range *p.enemies {
		if enemy == nil {
			continue
		}

		enemyRect := enemy.GetRect()
		c := rl.CheckCollisionCircleRec(v, float32(p.GetAttackRadius()), enemyRect)
		if c {
			enemy.ApplyDamage(p.GetDamage())
		}
	}

	go func() {
		<-p.attackCircleTimer.C
		p.mu.Lock()
		p.isAttacking = false
		p.mu.Unlock()
	}()
}

func (p *Player) GetAttackRadius() int {
	return p.attackRadius
}
