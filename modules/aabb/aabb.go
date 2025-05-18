package aabb

import (
	"meermookh/config"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollisionInfo struct {
	IsCollided bool
	IsStanding bool
	Entity     Drawable
	Side       string // "top", "bottom", "left", "right"
}

type Drawable interface {
	Draw(tex *rl.Texture2D)
	GetRect() *rl.Rectangle
}

func SimpleAABB(r1, r2 *rl.Rectangle) bool {
	if r1 == nil || r2 == nil {
		return false
	}

	return (r1.X < r2.X+r2.Width &&
		r1.X+r1.Width > r2.X &&
		r1.Y < r2.Y+r2.Height &&
		r1.Y+r1.Height > r2.Y)
}

func Check(plRect *rl.Rectangle, tiles *[]Drawable) CollisionInfo {
	filtered := split(plRect, tiles)

	if len(filtered) == 0 {
		return CollisionInfo{IsCollided: false}
	}

	for _, t := range filtered {
		tRect := t.GetRect()
		if rl.CheckCollisionRecs(*plRect, *tRect) {
			plBottom := plRect.Y + plRect.Height
			tTop := tRect.Y
			plRight := plRect.X + plRect.Width
			tLeft := tRect.X
			plLeft := plRect.X
			tRight := tRect.X + tRect.Width
			plTop := plRect.Y
			tBottom := tRect.Y + tRect.Height

			if plBottom <= tTop+5 && plBottom >= tTop-5 {
				return CollisionInfo{
					IsCollided: true,
					IsStanding: true,
					Entity:     t,
					Side:       "top",
				}
			}

			var side string
			if plBottom > tTop && plTop < tBottom {
				if plRight-tLeft < tRight-plLeft {
					side = "left"
				} else {
					side = "right"
				}
			} else {
				if plBottom-tTop < tBottom-plTop {
					side = "top"
				} else {
					side = "bottom"
				}
			}

			return CollisionInfo{
				IsCollided: true,
				IsStanding: false,
				Entity:     t,
				Side:       side,
			}
		}
	}

	return CollisionInfo{IsCollided: false}
}

// 1. Split window in 4 parts
// 2. Find out in which part player is
// 3. Check collision only with rects in such area
func split(plRect *rl.Rectangle, tiles *[]Drawable) (filtered []Drawable) {
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

	filtered = make([]Drawable, 0)
	for _, t := range *tiles {
		r := t.GetRect()
		if SimpleAABB(r, searchArea) {
			filtered = append(filtered, t)
		}
	}

	return filtered
}
