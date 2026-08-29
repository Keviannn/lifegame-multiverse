package world

import (
	"fmt"
	"sync"

	"github.com/Keviannn/lifegame-multiverse/internal/rules"
	"github.com/hajimehoshi/ebiten/v2"
)

// ###############    DEFINITIONS    ###############

const (
	Alive = true
	Dead = false

	// All the edge
	lefE = 0
	rigE = 1

	// Inside the edge [X, (...), Y]
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
	// World size data
	Width uint
	Heigh uint
	Size uint

	WorldImage *ebiten.Image
	view []byte

	// World corners
	botRC uint
	botLC uint
	topRC uint
	topLC uint

	// Where to move to check for neighbours
	cellNeighbours [8]int

	present, future uint8
	cells [2][]bool
	kind []uint8
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
		WorldImage: ebiten.NewImage(int(x), int(y)),
		view: make([]byte, x * y * 4),
		rules: *v,
		present: 0,
		future: 1,
	}

	w.topLC = 0
	w.topRC = w.Width - 1
	w.botLC = (w.Size - 1) - (w.Width - 1)
	w.botRC = w.Size - 1

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

// Returns cell state
func (w *World) GetCell(x, y uint) (bool, error) {
	pos := w.toAbs(x, y)

	if err := w.checkCoords(pos); err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return w.cells[w.present][pos], nil
}

// Returns cell state
func (w *World) getCellAbs(pos uint) (bool, error) {
	if err := w.checkCoords(pos); err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return w.cells[w.present][pos], nil
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

// Sets cell state in the future
func (w *World) setCellAbs(pos uint, state bool) error {
	if err := w.checkCoords(pos); err != nil {
		return fmt.Errorf("%w", err)
	}

	w.cells[w.future][pos] = state

	return nil
}

func (w *World) aliveNeighboursAbs(pos uint) (uint8, error) {
	err := w.checkCoords(pos)
	if err != nil {
		return 0, err
	}

	acc := 0

	for _, v := range chk[w.kind[pos]] {
		if w.cells[w.present][uint(int(pos) + w.cellNeighbours[v])] {
			acc += 1
		}
	}

	return uint8(acc), nil
}

// Simulates a new generation
func (w *World) NewGeneration() {
	w.drawWorld()
	for i := range w.Size {
		s, _ := w.getCellAbs(i)
		a, _ := w.aliveNeighboursAbs(i)

		n := w.rules.DecideFate(s, a)

		w.setCellAbs(i, n)
	}
	w.present = 1 - w.present
	w.future = 1 - w.future
}

// Simulates a new generation
// TODO: make dynamic routine assertion based in size
func (w *World) NewGenerationRoutines() {
	part := w.Size / 4
	var start uint = 0
	var finish uint = part

	var wg sync.WaitGroup
	for range 4 {
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
	w.drawWorld()
}

func (w *World) newGenPartial(from, until uint) {
	for i := from; i < until; i++ {
		s, _ := w.getCellAbs(i)
		a, _ := w.aliveNeighboursAbs(i)

		n := w.rules.DecideFate(s, a)

		w.setCellAbs(i, n)
	}
}

// Draws the world
func (w *World) drawWorld() {
	for i, v := range w.cells[w.present] {
		if v {
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
	w.WorldImage.WritePixels(w.view)
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
	isC := pos == w.topLC || pos == w.topRC || pos == w.botLC || pos == w.botRC
	isLE := rem == 0
	isRE := rem == w.Width - 1
	isTE := pos <= w.topRC
	isBE := pos >= w.botLC

	// Check from most common to least common celltype
	if  !isLE && !isRE && !isTE && !isBE {
		return inCell
	} else {
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
			case w.topLC:
				return topLC
			case w.topRC:
				return  topRC
			case w.botLC:
				return botLC
			default:
				return botRC
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
