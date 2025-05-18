package tile

import (
	"fmt"
	"meermookh/config"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	rect    *rl.Rectangle
	srcRect *rl.Rectangle // position of texture at the tileset
}

func New(pos rl.Vector2, idx int) Tile {
	rect := rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  config.BASE_TILE_SIZE,
		Height: config.BASE_TILE_SIZE,
	}

	srcRect := getSrcRect(idx)

	return Tile{
		rect:    &rect,
		srcRect: &srcRect,
	}
}

func (t *Tile) Draw(tex *rl.Texture2D) {
	fmt.Printf("tile.draw\n")
	if t.rect != nil && t.srcRect != nil {
		fmt.Printf("tile rects is not nil\n")
		rl.DrawTexturePro(
			*tex,
			*t.srcRect,
			*t.rect,
			rl.Vector2Zero(),
			0,
			rl.White,
		)
	}
}

func (t *Tile) GetRect() *rl.Rectangle {
	return t.rect
}

func getSrcRect(idx int) rl.Rectangle {
	idx = idx - 1
	row := idx / config.TILESET_WIDTH
	col := idx % config.TILESET_WIDTH

	src := rl.Rectangle{
		X:      float32(col * config.BASE_TILE_SIZE),
		Y:      float32(row * config.BASE_TILE_SIZE),
		Width:  config.BASE_TILE_SIZE,
		Height: config.BASE_TILE_SIZE,
	}

	fmt.Printf("src: %+v (idx: %d, row: %d, col: %d)\n", src, idx, row, col)

	return src
}
