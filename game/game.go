package game

import (
	"fmt"
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
	for pixelsFilled <= config.WINDOW_W {
		x := float32(pixelsFilled)
		y := float32(config.WINDOW_H - config.BASE_TILE_SIZE)

		fmt.Printf("x: %f\ty: %f\n", x, y)

		tile := tile.New(rl.Vector2{X: x, Y: y})
		pixelsFilled += config.BASE_TILE_SIZE
		tiles = append(tiles, tile)
	}

	fmt.Printf("Generated %d tiles.\n", len(tiles))

	return Game{
		tiles:  tiles,
		Player: player.New(),
	}
}

func (g *Game) Title(title string) {
	g.title = title
}

func (g *Game) Update() {
	aabb.Check(&g.Player, g.tiles)

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
