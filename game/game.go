package game

import (
	"meermookh/modules/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Player player.Player

	title string
}

func New() Game {
	return Game{
		Player: player.New(),
	}
}

func (g *Game) Title(title string) {
	g.title = title
}

func (g *Game) Update() {}

func (g *Game) Render() {
	rl.ClearBackground(rl.RayWhite)
	g.Player.Draw()
}

func (g *Game) Start() {
	rl.InitWindow(800, 450, g.title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		g.Update()
		g.Render()

		rl.EndDrawing()
	}
}
