package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g* Game) Update() error {
	// Update logical state
	return nil
}

func (g* Game) Draw(screen *ebiten.Image) {
	// Render the screen
}

func (g* Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Return the game logical screen size.
	// The screen is automatically scaled.
	return 320, 240
}

func main()  {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("LifeGame")

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
	
}
