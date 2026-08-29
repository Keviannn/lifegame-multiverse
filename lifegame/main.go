package main

import (
	"log"

	"github.com/Keviannn/lifegame-multiverse/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	World *world.World
	View *ebiten.Image
}

func (g* Game) Update() error {
	g.World.NewGenerationRoutines()
	return nil
}

func (g* Game) Draw(screen *ebiten.Image) {
	const cols = 2
	screen.DrawImage(g.World.WorldImage, nil)
}

func (g* Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main()  {
	
	u4, err := world.NewWorld(screenWidth, screenHeight, "B014/S2")
	if err != nil {
		log.Fatal(err)
	}

	u4.SetCell(91, 98, world.Alive)
	u4.SetCell(93, 99, world.Alive)

	u4.SetCell(90, 100, world.Alive)
	u4.SetCell(91, 100, world.Alive)
	u4.SetCell(94, 100, world.Alive)
	u4.SetCell(95, 100, world.Alive)
	u4.SetCell(96, 100, world.Alive)

	ebiten.SetWindowSize(screenWidth * 4, screenHeight * 4)
	ebiten.SetWindowTitle("LifeGame")

	game := &Game {
		World: u4,
		View: ebiten.NewImage(screenWidth, screenHeight),
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
