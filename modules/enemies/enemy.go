package enemies

import (
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/interfaces"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EnemyAction string

const (
	EnemyPatrolAction EnemyAction = "patrol"
	EnemyChaseAction  EnemyAction = "chase"
	EnemyAttackAction EnemyAction = "attack"
)

type Enemy struct {
	rect       rl.Rectangle
	hp         float32
	isStanding bool
	speed      float32
	gravity    float32
	damage     int

	detectionAreaRadius float32
	attackAreaRadius    float32

	color rl.Color
	state EnemyAction

	plRect *rl.Rectangle
	player interfaces.Player

	isAttacking    bool
	attackTimer    *time.Timer
	attackCooldown time.Duration

	mu sync.Mutex
}

func New(pos rl.Vector2) Enemy {
	return Enemy{
		isStanding: false,
		// TODO
		// State is hardcoded for debug purposes
		state:               EnemyChaseAction,
		speed:               2,
		detectionAreaRadius: 128,
		attackCooldown:      750 * time.Millisecond,
		attackAreaRadius:    64,
		damage:              10,
		hp:                  100,
		gravity:             5,
		color:               rl.Red,
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

	e.manageState()

	coll := aabb.Check(&e.rect, tiles)
	if !coll.IsStanding {
		e.rect.Y += e.gravity
	}
}

func (e *Enemy) AttachPlayerRectPtr(r *rl.Rectangle) {
	e.plRect = r
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

func (e *Enemy) SetPlayer(p interfaces.Player) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.player = p
}

func (e *Enemy) ChaseState() {
	c := rl.Vector2{
		X: e.rect.X,
		Y: e.rect.Y,
	}

	rl.DrawCircleLines(e.rect.ToInt32().X, e.rect.ToInt32().Y, e.detectionAreaRadius, rl.Orange)
	if rl.CheckCollisionCircleRec(c, e.detectionAreaRadius, *e.plRect) {

		plPos := rl.Vector2{
			X: e.plRect.X,
			Y: e.plRect.Y,
		}

		enPos := rl.Vector2{
			X: e.rect.X,
			Y: e.rect.Y,
		}

		direction := rl.Vector2Subtract(plPos, enPos)
		direction = rl.Vector2Normalize(direction)
		velocity := rl.Vector2Scale(direction, e.speed)

		e.rect.X += velocity.X
		e.rect.Y += velocity.Y

		if rl.CheckCollisionCircleRec(c, e.attackAreaRadius, *e.plRect) {
			rl.DrawCircleLines(e.rect.ToInt32().X, e.rect.ToInt32().Y, e.attackAreaRadius, rl.Pink)

			if !e.isAttacking && e.player != nil {
				e.isAttacking = true
				e.attackTimer = time.NewTimer(e.attackCooldown)
				e.player.DealDamage(e.damage)

				go func() {
					<-e.attackTimer.C
					e.isAttacking = false
				}()
			}
		}
	}
}

func (e *Enemy) PatrolState() {
	e.mu.Lock()
	defer e.mu.Unlock()

}

func (e *Enemy) AttackState() {
	e.mu.Lock()
	defer e.mu.Unlock()

}

// State
func (e *Enemy) manageState() {
	switch e.state {
	case EnemyPatrolAction:
		e.PatrolState()
	case EnemyChaseAction:
		e.ChaseState()
	case EnemyAttackAction:
		e.AttackState()
	}
}
