package game

import (
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/player"
	"meermookh/modules/tile"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Player player.Player

	tiles []tile.Tile
	title string
}

func New() Game {
	tiles := make([]tile.Tile, 0)
	pixelsFilled := 0
	for pixelsFilled <= config.WINDOW_W/2 {
		x := float32(pixelsFilled)
		y := float32(config.WINDOW_H - config.BASE_TILE_SIZE)

		tile := tile.New(rl.Vector2{X: x, Y: y})
		pixelsFilled += config.BASE_TILE_SIZE
		tiles = append(tiles, tile)
	}

	return Game{
		tiles:  tiles,
		Player: player.New(),
	}
}

func (g *Game) Title(title string) {
	g.title = title
}

func (g *Game) Update() {
	plRect := g.Player.GetRect()
	info := aabb.Check(&plRect, g.tiles)
	if info.IsCollided {
		g.Player.HandleCollision(info)
	} else {
		g.Player.ResetCollision()
	}

	g.Player.Update()
}

func (g *Game) Render() {
	rl.ClearBackground(rl.RayWhite)

	for _, tile := range g.tiles {
		tile.Draw()
	}

	g.Player.Draw()
}

func (g *Game) Start() {
	rl.InitWindow(config.WINDOW_W, config.WINDOW_H, g.title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		g.Update()
		g.Render()

		rl.EndDrawing()
	}
}
