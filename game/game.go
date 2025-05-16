package game

import (
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/enemies"
	"meermookh/modules/player"
	"meermookh/modules/tile"

	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Player player.Player

	tiles           []tile.Tile
	title           string
	enemies         []*enemies.Enemy
	enemiesToDelete []int
}

func New() Game {
	tiles := make([]tile.Tile, 0)
	enemiesSlice := make([]*enemies.Enemy, 0)
	pixelsFilled := 0
	for pixelsFilled <= config.WINDOW_W/2 {
		x := float32(pixelsFilled)
		y := float32(config.WINDOW_H - config.BASE_TILE_SIZE)

		tile := tile.New(rl.Vector2{X: x, Y: y})
		pixelsFilled += config.BASE_TILE_SIZE
		tiles = append(tiles, tile)
	}

	for i := range 10 {
		pos := rl.Vector2{
			X: float32(250 + i*100),
			Y: 500,
		}

		enemy := enemies.New(pos)
		enemiesSlice = append(enemiesSlice, &enemy)
	}

	return Game{
		enemiesToDelete: make([]int, 0),
		tiles:           tiles,
		Player:          player.New(rl.Vector2{X: 100, Y: 700}),
		enemies:         enemiesSlice,
	}
}

func (g *Game) Title(title string) {
	g.title = title
}

func (g *Game) Update() {
	// move player collision handling to Player.go
	plRect := g.Player.GetRect()
	info := aabb.Check(&plRect, &g.tiles)
	if info.IsCollided {
		g.Player.HandleCollision(info)
	} else {
		g.Player.ResetCollision()
	}

	g.ManageEnemies()

	for i := range g.enemies {
		if g.enemies[i] != nil {
			go g.enemies[i].Update(&g.tiles)
		}
	}

	g.Player.Update()
}

func (g *Game) ManageEnemies() {
	for i, enemy := range g.enemies {
		if enemy != nil {
			rect := enemy.GetRect()
			if rect.X >= config.WINDOW_W || rect.Y >= config.WINDOW_H {
				g.enemiesToDelete = append(g.enemiesToDelete, i)
			}
		}
	}

	for i := len(g.enemiesToDelete) - 1; i >= 0; i-- {
		if g.enemiesToDelete[i] < len(g.enemies) {
			g.enemies = slices.Delete(g.enemies, g.enemiesToDelete[i], g.enemiesToDelete[i]+1)
		}
	}

	g.enemiesToDelete = g.enemiesToDelete[:0]
}

func (g *Game) Render() {
	rl.ClearBackground(rl.RayWhite)

	for _, tile := range g.tiles {
		tile.Draw()
	}

	for _, enemy := range g.enemies {
		enemy.Draw()
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

func (g *Game) GetTiles() *[]tile.Tile {
	return &g.tiles
}
