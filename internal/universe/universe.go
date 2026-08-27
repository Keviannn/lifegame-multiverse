package universe 

import (
	"fmt"

	"github.com/Keviannn/lifegame-multiverse/internal/rules"
	"github.com/Keviannn/lifegame-multiverse/internal/world"
)

type Universe struct {
	World *world.World
	rules *rules.Rules
}

// Creates a new universe
func NewUniverse (x, y uint, r string) (*Universe, error) {
	v, err := rules.NewRules(r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &Universe{
		World: world.NewWorld(x, y),
		rules: v,
	}, nil
}

// Simulates a new generation
func (u *Universe) NewGeneration() {
	otherWorld := world.NewWorld(u.World.Width, u.World.Heigh)
	for i := range u.World.Size {
		s, _ := u.World.GetCellAbs(i)
		a, _ := u.World.AliveNeighboursAbs(i)

		n := u.rules.DecideFate(s, a)
		
		otherWorld.SetCellAbs(i, n)
	}

	u.World.SetWorldState(otherWorld.GetWorldState())
}

// Draws the universe to a byte slice
func (u *Universe) DrawUniverse(pixels []byte) {
	for i, v := range u.World.GetWorldState() {
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
