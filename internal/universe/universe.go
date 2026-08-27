package universe 

import (
	"fmt"

	"github.com/Keviannn/lifegame-multiverse/internal/rules"
	"github.com/Keviannn/lifegame-multiverse/internal/world"
)

type Universe struct {
	world *world.World
	rules *rules.Rules
}

// Creates a new universe
func NewUniverse (x, y uint, r string) (*Universe, error) {
	v, err := rules.NewRules(r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &Universe{
		world: world.NewWorld(x, y),
		rules: v,
	}, nil
}

// Simulates a new generation
func (u *Universe) NewGeneration() {
	otherWorld := world.NewWorld(u.world.Width, u.world.Heigh)
	for i := range u.world.Size {
		s, _ := u.world.GetCellAbs(i)
		a, _ := u.world.AliveNeighboursAbs(i)

		n := u.rules.DecideFate(s, a)

		
		otherWorld.SetCellAbs(i, n)
	}

	u.world.SetCells(otherWorld.GetCells())
}

// Draws the universe to a byte slice
func (u *Universe) DrawUniverse(pixels []byte) {
	for i, v := range u.world.GetCells() {
		if v {
			pixels[4*i] = 0xff
			pixels[4*i+1] = 0xff
			pixels[4*i+2] = 0xff
			pixels[4*i+3] = 0xff
		} else {
			pixels[4*i] = 0
			pixels[4*i+1] = 0
			pixels[4*i+2] = 0
			pixels[4*i+3] = 0
		}
	}
}
