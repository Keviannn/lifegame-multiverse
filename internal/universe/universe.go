package universe

import (
	"fmt"
	"sync"

	"github.com/Keviannn/lifegame-multiverse/internal/rules"
	"github.com/Keviannn/lifegame-multiverse/internal/world"
)

type Universe struct {
	Worlds *[2]world.World
	rules *rules.Rules
	status uint
}

// Creates a new universe
func NewUniverse (x, y uint, r string) (*Universe, error) {
	v, err := rules.NewRules(r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	w := [2]world.World {*world.NewWorld(x, y), *world.NewWorld(x, y)}

	return &Universe {
		Worlds: &w,
		rules: v,
		status: 0,
	}, nil
}

// Simulates a new generation
func (u *Universe) NewGeneration() {
	org, dest := u.status, 1 - u.status
	for i := range u.Worlds[0].Size {
		s, _ := u.Worlds[org].GetCellAbs(i)
		a, _ := u.Worlds[org].AliveNeighboursAbs(i)

		n := u.rules.DecideFate(s, a)
		
		u.Worlds[dest].SetCellAbs(i, n)
	}

	u.status = dest
}

func (u *Universe) NewGenerationRoutines() {
	org, dest := u.status, 1 - u.status
	part := u.Worlds[0].Size / 4
	var start uint = 0
	var finish uint = part

	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func(from, until, org, dest uint) {
			defer wg.Done()
			u.newGenPartial(from, until, org, dest)
		}(start, finish, org, dest)
		start = finish
		finish += part 
	}
	wg.Wait()
	u.status = dest
}

func (u *Universe) newGenPartial(from, until, org, dest uint) {
	for i := from; i < until; i++ {
		s, _ := u.Worlds[org].GetCellAbs(i)
		a, _ := u.Worlds[org].AliveNeighboursAbs(i)

		n := u.rules.DecideFate(s, a)
		
		u.Worlds[dest].SetCellAbs(i, n)
	}
}

// Draws the universe to a byte slice
func (u *Universe) DrawUniverse(pixels []byte) {
	for i, v := range u.Worlds[0].GetWorldState() {
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
