package game

import (
	"meermookh/config"
	"meermookh/modules/aabb"
	"meermookh/modules/enemies"
	"meermookh/modules/events"
	"meermookh/modules/player"
	"meermookh/modules/tile"
	"sync"

	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	mu sync.RWMutex

	shouldRun       bool
	player          *player.Player
	tiles           []tile.Tile
	title           string
	enemies         []*enemies.Enemy
	enemiesToDelete []int

	currentScreen config.ScreenType
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

	pl := player.New(rl.Vector2{X: 100, Y: 700}, &enemiesSlice)

	return Game{
		enemiesToDelete: make([]int, 0),
		tiles:           tiles,
		player:          &pl,
		enemies:         enemiesSlice,
	}
}

func (g *Game) Update() {
	g.mu.Lock()
	g.ManageEnemies()

	drawables := make([]aabb.Drawable, len(g.tiles))
	for i := range g.tiles {
		drawables[i] = &g.tiles[i]
	}

	g.player.Update(&drawables)
	if g.player.GetHP() <= 0 {
		g.currentScreen = config.DeadScreen
	}
	g.mu.Unlock()

	var wg sync.WaitGroup
	g.mu.RLock()
	enemiesCopy := make([]*enemies.Enemy, len(g.enemies))
	copy(enemiesCopy, g.enemies)
	g.mu.RUnlock()

	for _, enemy := range enemiesCopy {
		if enemy != nil {
			wg.Add(1)
			go func(e *enemies.Enemy) {
				defer wg.Done()
				e.Update(&drawables)
			}(enemy)
		}
	}
	wg.Wait()
}

func (g *Game) ManageEnemies() {
	for i, enemy := range g.enemies {
		if enemy != nil {
			// all checks for enemies
			rect := enemy.GetRect()
			if rect.X >= config.WINDOW_W || rect.Y >= config.WINDOW_H {
				g.enemiesToDelete = append(g.enemiesToDelete, i)
			}

			if enemy.GetHP() <= 0 {
				events.PlayerKilledEnemy(g.player)
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
	for _, tile := range g.tiles {
		tile.Draw()
	}

	for _, enemy := range g.enemies {
		if enemy != nil {
			enemy.Draw()
		}
	}

	g.player.Draw()
}

func (g *Game) Start() {
	rl.InitWindow(config.WINDOW_W, config.WINDOW_H, "meer mookh")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	g.shouldRun = true
	g.currentScreen = config.StartScreen
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		switch g.currentScreen {
		case config.StartScreen:
			g.DrawStartScreen()

		case config.GameScreen:
			g.DrawGameScreen()

		case config.DeadScreen:
			g.DrawDeadScreen()
		}

		if !g.shouldRun {
			break
		}

		rl.EndDrawing()
	}
}

func (g *Game) GetTiles() *[]tile.Tile {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return &g.tiles
}

func (g *Game) DrawGameScreen() {
	g.Update()
	g.Render()
}

func (g *Game) DrawDeadScreen() {
	rl.DrawText("TI UMER!!!! DOLBOEB", 600, 400, 32, rl.Black)
	rl.DrawText("press space to quit", 600, 500, 32, rl.Black)

	if rl.IsKeyPressed(rl.KeySpace) {
		g.shouldRun = false
	}
}

func (g *Game) DrawStartScreen() {
	rl.DrawText("press space to play", 600, 600, 32, rl.Black)

	if rl.IsKeyPressed(rl.KeySpace) {
		g.currentScreen = config.GameScreen
	}
}
