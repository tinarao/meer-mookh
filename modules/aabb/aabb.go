package aabb

import (
	"fmt"
	"meermookh/config"
	"meermookh/modules/player"
	"meermookh/modules/tile"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SimpleAABB(r1, r2 *rl.Rectangle) bool {
	if r1 == nil || r2 == nil {
		return false
	}

	return (r1.X < r2.X+r2.Width &&
		r1.X+r1.Width > r2.X &&
		r1.Y < r2.Y+r2.Height &&
		r1.Y+r1.Height > r2.Y)
}

func Check(pl *player.Player, tiles []tile.Tile) {
	filtered := split(pl, tiles)
	fmt.Printf("%d tiles to check\n", len(filtered))

	if len(filtered) == 0 {
		return
	}
}

// 1. Split window in 4 parts
// 2. Find out in which part player is
// 3. Check collision only with rects in such area
func split(pl *player.Player, tiles []tile.Tile) (filtered []tile.Tile) {
	plRect := pl.GetRect()
	wCenter := config.WINDOW_W / 2
	hCenter := config.WINDOW_H / 2

	isLeft := plRect.X < float32(wCenter)
	isTop := plRect.Y < float32(hCenter)

	searchArea := &rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(wCenter),
		Height: float32(hCenter),
	}

	if !isTop {
		searchArea.Y = float32(hCenter)
		searchArea.Height = float32(config.WINDOW_H - hCenter)
	}

	if !isLeft {
		searchArea.X = float32(wCenter)
		searchArea.Width = float32(config.WINDOW_W - wCenter)
	}

	// Now, when i have the search area
	// i can filter tiles and get such tiles
	// that locates in this area
	// then check collisions with those
	filtered = make([]tile.Tile, 0)

	for _, t := range tiles {
		r := t.GetRect()
		if SimpleAABB(r, searchArea) {
			filtered = append(filtered, t)
		}
	}

	return filtered
}
