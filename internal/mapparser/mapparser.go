package mapparser

import (
	"encoding/xml"
	"fmt"
	"meermookh/modules/tile"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TMXMap struct {
	XMLName  xml.Name     `xml:"map"`
	Width    int          `xml:"width,attr"`
	Height   int          `xml:"height,attr"`
	Tilesets []TMXTileset `xml:"tileset"`
	Layers   []TMXLayer   `xml:"layer"`
}

type TMXTileset struct {
	FirstGID int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

type TMXLayer struct {
	Name   string `xml:"name,attr"`
	Data   string `xml:"data"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type Layer [30][60]string

type Tilemap struct {
	Layers map[int]Layer
	Tiles  []*tile.Tile
}

const DEFAULT_FLOOR_TILE_TYPE = "3"

func LoadMap(name string) *Tilemap {
	path := filepath.Join(".", "tmx", name)
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	tmxMap := &TMXMap{}
	if err := xml.Unmarshal(file, tmxMap); err != nil {
		panic(err)
	}

	tilemap := &Tilemap{
		Layers: make(map[int]Layer),
		Tiles:  make([]*tile.Tile, 0),
	}

	for i, layerData := range tmxMap.Layers {
		var layer Layer
		cleanData := strings.ReplaceAll(strings.TrimSpace(layerData.Data), "\n", "")
		splitted := strings.Split(cleanData, ",")

		if len(splitted) < 30*60 {
			panic("Not enough tile data in layer")
		}

		for y := range 30 {
			for x := range 60 {
				idx := y*60 + x
				if idx < len(splitted) {
					layer[y][x] = splitted[idx]
				}
			}
		}

		tilemap.Layers[i] = layer
	}

	tilemap.generateTiles()
	return tilemap
}

func (m *Tilemap) generateTiles() {
	for _, layer := range m.Layers {
		for y := range 30 {
			for x := range 60 {
				tileID := layer[y][x]
				if tileID != "0" {
					pos := rl.Vector2{
						X: float32(x * 32),
						Y: float32(y * 32),
					}

					idx, err := strconv.Atoi(tileID)
					if err != nil {
						fmt.Printf("Failed to convert tile ID %s to int: %v\n", tileID, err)
						continue
					}

					tile := tile.New(pos, idx)
					m.Tiles = append(m.Tiles, &tile)
				}
			}
		}
	}
}
