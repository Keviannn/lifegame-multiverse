package main

import (
	"log"

	"github.com/Keviannn/lifegame-multiverse/internal/universe"
	"github.com/Keviannn/lifegame-multiverse/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480 
)

type Game struct {
	universe *universe.Universe
	pixels []byte
}

func (g* Game) Update() error {
	g.universe.NewGeneration()
	return nil
}

func (g* Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, 4 * screenWidth * screenHeight)
	}
	g.universe.DrawUniverse(g.pixels)
	screen.WritePixels(g.pixels)
}

func (g* Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main()  {
	u, err := universe.NewUniverse(screenWidth, screenHeight, "B3/S23")
	if err != nil {
		log.Fatal(err)
	}

	
	u.World.SetCell(161, 118, world.Alive)
	u.World.SetCell(163, 119, world.Alive)

	u.World.SetCell(160, 120, world.Alive)
	u.World.SetCell(161, 120, world.Alive)
	u.World.SetCell(164, 120, world.Alive)
	u.World.SetCell(165, 120, world.Alive)
	u.World.SetCell(166, 120, world.Alive)

	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("LifeGame")

	game := &Game {
		universe: u,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
