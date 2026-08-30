package main

import (
	"fmt"
	"log"

	"github.com/Keviannn/lifegame-multiverse/internal/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 2560
	screenHeight = 1440
)

type Game struct {
	World *world.World
	Image *ebiten.Image
}

func (g* Game) Update() error {
	g.World.NewGenerationRoutines(16)
	return nil
}

func (g* Game) Draw(screen *ebiten.Image) {
	g.Image.WritePixels(g.World.DrawWorld(16))
	screen.DrawImage(g.Image, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Current TPS: %.0f\nCurrent FPS: %.0f\nWorld Width: %d\nWorld Height: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), screenWidth, screenHeight))
}

func (g* Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main()  {
	fullScreenWidth, fullScreenHeight := ebiten.Monitor().Size()
	u4, err := world.NewWorld(screenWidth, screenHeight, "B014/S2")
	if err != nil {
		log.Fatal(err)
	}

	u4.Populate(0.02)

	ebiten.SetWindowSize(fullScreenWidth, fullScreenHeight)
	ebiten.SetWindowTitle("LifeGame")
	ebiten.SetFullscreen(true)

	game := &Game {
		World: u4,
		Image: ebiten.NewImage(screenWidth, screenHeight),
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
