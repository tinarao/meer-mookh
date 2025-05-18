package config

const (
	WINDOW_W       = 1600
	WINDOW_H       = 960
	BASE_TILE_SIZE = 32
	PLAYER_WIDTH   = 32
	PLAYER_HEIGHT  = 32

	TILESET_WIDTH = 11 // 32x32 rects in row
)

type ScreenType string

const (
	StartScreen ScreenType = "start"
	GameScreen  ScreenType = "game"
	DeadScreen  ScreenType = "dead"
)
