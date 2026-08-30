package main

import (
	"fmt"
	"log"

	"github.com/Keviannn/lifegame-multiverse/internal/world"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 2560
	screenHeight = 1920
)

type Game struct {
	World *world.World
	View *ebiten.Image
}

func (g* Game) Update() error {
	g.World.NewGenerationRoutines(16)
	return nil
}

func (g* Game) Draw(screen *ebiten.Image) {
	tps := fmt.Sprintf("Current TPS: %.0f\nCurrent FPS: %.0f", ebiten.ActualTPS(), ebiten.ActualFPS())
	screen.WritePixels(g.World.DrawWorld(16))
	ebitenutil.DebugPrint(screen, tps)
}

func (g* Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main()  {
	
	u4, err := world.NewWorld(screenWidth, screenHeight, "B014/S2")
	if err != nil {
		log.Fatal(err)
	}

	u4.Populate(0.05)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("LifeGame")

	game := &Game {
		World: u4,
		View: ebiten.NewImage(screenWidth, screenHeight),
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
