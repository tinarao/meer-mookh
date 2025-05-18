package interfaces

import rl "github.com/gen2brain/raylib-go/raylib"

type Player interface {
	DealDamage(amount int)
	GetRect() rl.Rectangle
}
