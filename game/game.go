package game

import (
	"fmt"
	"meermookh/config"
	"meermookh/internal/mapparser"
	"meermookh/modules/aabb"
	"meermookh/modules/enemies"
	"meermookh/modules/events"
	"meermookh/modules/player"
	"meermookh/modules/tile"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	mu sync.RWMutex

	shouldRun       bool
	player          *player.Player
	title           string
	enemies         []*enemies.Enemy
	enemiesToDelete []int

	tilemap        *mapparser.Tilemap
	loadedTextures map[string]*rl.Texture2D

	currentScreen config.ScreenType
}

func New() Game {
	enemiesSlice := make([]*enemies.Enemy, 0)

	for i := range 1 {
		pos := rl.Vector2{
			X: float32(250 + i*100),
			Y: 500,
		}

		enemy := enemies.New(pos)
		enemiesSlice = append(enemiesSlice, &enemy)
	}

	tilemap := mapparser.LoadMap("main_map.tmx")

	pl := player.New(rl.Vector2{X: 100, Y: 700}, &enemiesSlice)

	for _, e := range enemiesSlice {
		e.AttachPlayerRectPtr(pl.GetRectPtr())
		e.SetPlayer(&pl)
	}

	return Game{
		enemiesToDelete: make([]int, 0),
		tilemap:         tilemap,
		player:          &pl,
		enemies:         enemiesSlice,
		loadedTextures:  make(map[string]*rl.Texture2D),
	}
}

func (g *Game) Update() {
	g.mu.Lock()
	g.ManageEnemies()

	drawables := make([]aabb.Drawable, len(g.tilemap.Tiles))
	for i := range g.tilemap.Tiles {
		drawables[i] = g.tilemap.Tiles[i]
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
	g.DrawBackground()

	for _, tile := range g.tilemap.Tiles {
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

	g.loadTextures()
	defer g.unloadTextures()

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

func (g *Game) GetTiles() *[]*tile.Tile {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return &g.tilemap.Tiles
}

func (g *Game) DrawBackground() {
	g.mu.RLock()
	defer g.mu.RUnlock()

	tex := g.loadedTextures["game-bg-sky"]
	if tex == nil {
		panic("\"game-bg-sky\" texture is not present in loaded textures")
	}

	rl.DrawTextureEx(*tex, rl.Vector2Zero(), 0, 0.90, rl.Gray)
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

func (g *Game) loadTextures() {
	path := filepath.Join(".", "assets")
	dir, err := os.Open(path)
	if err != nil {
		e := fmt.Sprintf("failed to os.Open: %s", err.Error())
		panic(e)
	}

	files, err := dir.ReadDir(-1)
	if err != nil {
		e := fmt.Sprintf("failed to dir.ReadDir: %s", err.Error())
		panic(e)
	}

	allowedFormats := []string{
		"png",
		"jpg",
	}

	for _, file := range files {
		splitted := strings.Split(file.Name(), ".")
		if len(splitted) != 2 {
			e := fmt.Sprintf("Invalid file found at /assets/%s\n", file.Name())
			panic(e)
		}

		if !slices.Contains(allowedFormats, splitted[1]) {
			e := fmt.Sprintf("Invalid file found at /assets/%s\n", file.Name())
			panic(e)
		}

		filename := splitted[0]
		fullpath := filepath.Join("assets", file.Name())
		tex := rl.LoadTexture(fullpath)
		g.loadedTextures[filename] = &tex
	}

	fmt.Printf("Loaded %d textures.\n", len(g.loadedTextures))

}

func (g *Game) unloadTextures() {
	for _, t := range g.loadedTextures {
		rl.UnloadTexture(*t)
	}
}
