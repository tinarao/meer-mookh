package config

const (
	WINDOW_W       = 1600
	WINDOW_H       = 900
	BASE_TILE_SIZE = 32
	PLAYER_WIDTH   = 32
	PLAYER_HEIGHT  = 32
)

type ScreenType string

const (
	StartScreen ScreenType = "start"
	GameScreen  ScreenType = "game"
	DeadScreen  ScreenType = "dead"
)
