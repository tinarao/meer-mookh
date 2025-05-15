package player

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	Rect rl.Rectangle
	Size uint
}

func New() Player {
	return Player{
		Size: 32,
		Rect: rl.Rectangle{
			X:      100,
			Y:      100,
			Width:  32,
			Height: 32,
		},
	}
}

func (p *Player) Draw() {
	origin := rl.Vector2{
		X: float32(p.Size) / 2,
		Y: float32(p.Size) / 2,
	}

	rl.DrawRectanglePro(p.Rect, origin, 0, rl.Red)
}
