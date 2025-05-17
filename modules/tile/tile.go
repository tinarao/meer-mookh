package tile

import (
	"meermookh/config"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	rect rl.Rectangle
}

func New(pos rl.Vector2) Tile {
	return Tile{
		rect: rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  config.BASE_TILE_SIZE,
			Height: config.BASE_TILE_SIZE,
		},
	}
}

func (t *Tile) Draw() {
	rl.DrawRectangleRec(t.rect, rl.Green)
}

func (t *Tile) GetRect() *rl.Rectangle {
	return &t.rect
}
