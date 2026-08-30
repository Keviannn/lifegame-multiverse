package world

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/Keviannn/lifegame-multiverse/internal/rules"
)

// ###############    DEFINITIONS    ###############

const (
	Alive = true
	Dead = false

	lefE = 0
	rigE = 1
	topE = 2
	botE = 3
	topLC = 4
	topRC = 5
	botLC = 6
	botRC = 7
	inCell = 8
)

// What neighbours to visit for every cell position
var chk = [9][]uint8 {
	{1, 2, 4, 6, 7},
	{0, 1, 3, 5, 6},
	{3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4},
	{4, 6, 7},
	{3, 5, 6},
	{1, 2, 4},
	{0, 1, 3},
	{0, 1, 2, 3, 4, 5, 6, 7},
}

// This is where they live
// Must be initialized after specifying Width and Heigh
type World struct {
	Width uint
	Heigh uint
	Size uint

	// Ebiten view
	view []byte

	// Two generations
	present, future uint8
	cells [2][]bool

	// What neighbours to visit
	cellNeighbours [8]int

	// Relative position of all cells
	kind []uint8

	// World rules
	rules rules.Rules
}


// ###############    PUBLIC METHODS    ###############

// Creates the slice of cells based on worlds width and height
// and gives values to important positions in the world
func NewWorld(x, y uint, r string) (*World, error) {
	v, err := rules.NewRules(r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	w := &World {
		Width: x,
		Heigh: y,
		Size: x * y,
		view: make([]byte, x * y * 4),
		rules: *v,
		present: 0,
		future: 1,
	}

	w.cellNeighbours = [8]int {
		int(-w.Width -1), int(-w.Width), int(-w.Width + 1),
		-1,		 			    		   				1 ,
		int(w.Width - 1), int(w.Width),  int(w.Width +  1),
	}

	for i := range 2 {
		w.cells[i] = make([]bool, w.Size)
	}

	w.kind = make([]uint8, w.Size)
	w.checkAllCases()

	return w, nil
}

// Populate the world with selected density
func (w *World) Populate (density float64) {
	for i := range w.Size {
		if rand.Float64() < density {
			w.cells[0][i] = true
		}
	}
}

// Returns cell state
func (w *World) GetCell(x, y uint) (bool, error) {
	pos := w.toAbs(x, y)

	if err := w.checkCoords(pos); err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return w.cells[w.present][pos], nil
}

// Returns cell state in the present for absolute position
func (w *World) getCellAbs(pos uint) bool {
	return w.cells[w.present][pos]
}

// Sets cell state
func (w *World) SetCell(x, y uint, state bool) error {
	pos := w.toAbs(x, y)

	if err := w.checkCoords(pos); err != nil {
		return fmt.Errorf("%w", err)
	}

	w.cells[w.present][pos] = state

	return nil
}

// Sets cell state in the future for absolute position
func (w *World) setCellAbs(pos uint, state bool) {
	w.cells[w.future][pos] = state
}

// Counts all alive neighbours of a cell
func (w *World) aliveNeighboursAbs(pos uint) uint8 {
	acc := 0

	for _, v := range chk[w.kind[pos]] {
		if w.cells[w.present][uint(int(pos) + w.cellNeighbours[v])] {
			acc += 1
		}
	}

	return uint8(acc)
}

// Simulates a new generation
func (w *World) NewGeneration() {
	for i := range w.Size {
		s := w.getCellAbs(i)
		a := w.aliveNeighboursAbs(i)

		n := w.rules.DecideFate(s, a)

		w.setCellAbs(i, n)
	}
	w.present = 1 - w.present
	w.future = 1 - w.future
}

// Simulates a new generation
// TODO: make dynamic routine assertion based in size
func (w *World) NewGenerationRoutines(num uint) {
	part := w.Size / num
	var start uint = 0
	var finish uint = part

	var wg sync.WaitGroup
	for range num {
		wg.Add(1)
		go func(from, until uint) {
			defer wg.Done()
			w.newGenPartial(from, until)
		}(start, finish)
		start = finish
		finish += part
	}
	wg.Wait()
	w.present = 1 - w.present
	w.future = 1 - w.future
}

func (w *World) newGenPartial(from, until uint) {
	for i := from; i < until; i++ {
		s := w.getCellAbs(i)
		a := w.aliveNeighboursAbs(i)

		n := w.rules.DecideFate(s, a)

		w.setCellAbs(i, n)
	}
}

// Draws the world and sets its image
func (w *World) DrawWorld(num uint) []byte {
	part := w.Size / num
	var start uint = 0
	var finish uint = part

	var wg sync.WaitGroup
	for range num {
		wg.Add(1)
		go func(from, until uint) {
			defer wg.Done()
			w.drawPartial(from, until)
		}(start, finish)
		start = finish
		finish += part
	}
	wg.Wait()
	return w.view
}

func (w *World) drawPartial (from, until uint) {
	for i := from; i < until; i++ {
		if w.cells[w.present][i] {
			w.view[4*i] = 0xff
			w.view[4*i+1] = 0xff
			w.view[4*i+2] = 0xff
			w.view[4*i+3] = 0xff
		} else {
			w.view[4*i] = 0
			w.view[4*i+1] = 0
			w.view[4*i+2] = 0
			w.view[4*i+3] = 0
		}
	}
}

// ###############    PRIVATE METHODS    ###############

// Converts (x, y) coordinates to absolute position 
func (w *World) toAbs(x, y uint) uint {
	return y * w.Width + x
}

// Checks if coordinates are inside the world
func (w *World) checkCoords(pos uint) error {
	if pos >= w.Size {
		return fmt.Errorf("Cell coordinates out of bounds for a %dx%d world", w.Width, w.Heigh)
	}

	return nil
}

// Returns relative position of cell like top left corner
// Does not check for error as its called always after checking
func (w *World) checkCase(pos uint) uint8 {
	rem := pos % w.Width
	isLE := rem == 0
	isRE := rem == w.Width - 1
	isTE := pos <= w.Width - 1
	isBE := pos >= (w.Size - 1) - (w.Width - 1)

	// Check from most common to least common celltype
	if  !isLE && !isRE && !isTE && !isBE {
		return inCell
	} else {
		isC := pos == 0 || pos == w.Width - 1 || pos == (w.Size - 1) - (w.Width - 1) || pos == w.Size - 1
		if isLE && !isC {
			return lefE
		} else if isRE && !isC {
			return rigE
		} else if isTE && !isC {
			return topE
		} else if isBE && !isC {
			return botE
		} else {
			switch pos {
			case 0:
				return topLC
			case w.Width - 1:
				return  topRC
			case w.Size - 1:
				return botRC
			default:
				return botLC
			}
		}
	}
}

// Precalculates all relative positions for all cells using goroutines
// TODO: make dynamic routine assertion based in size
func (w *World) checkAllCases () {
	part := w.Size / 4
	var start uint = 0
	var finish uint = part
	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func(from, until uint) {
			defer wg.Done()
			for i := from; i < until; i++ {
				w.kind[i] = w.checkCase(uint(i))
			}
		}(start, finish)
		start = finish
		finish += part
	}
	wg.Wait()
}
