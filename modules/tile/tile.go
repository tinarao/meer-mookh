package tile

import (
	"meermookh/config"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	rect    rl.Rectangle
	isSolid bool // should player collide with this tile
}

func New(pos rl.Vector2) Tile {
	return Tile{
		isSolid: true,
		rect: rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  config.BASE_TILE_SIZE,
			Height: config.BASE_TILE_SIZE,
		},
	}
}

func (t *Tile) Draw() {
	origin := rl.Vector2{
		X: float32(config.BASE_TILE_SIZE) / 2,
		Y: float32(config.BASE_TILE_SIZE) / 2,
	}

	rl.DrawRectanglePro(t.rect, origin, 0, rl.Green)
}

func (t *Tile) GetRect() *rl.Rectangle {
	return &t.rect
}
