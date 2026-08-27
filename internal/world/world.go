package world

import (
	"fmt"
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
	botLC = 5
	topRC = 6
	botRC = 7
	inCell = 8
)

type Coords struct {
	x int
	y int
}

// Full check
var chkF = []uint8 {0, 1, 2, 3, 4, 5, 6, 7}

// Edge checks
var chkTE = []uint8 {0, 1, 2, 3, 4}
var chkBE = []uint8 {3, 4, 5, 6, 7}
var chkLE = []uint8 {0, 1, 3, 5, 6}
var chkRE = []uint8 {1, 2, 4, 6, 7} 

// Corner checks
var chkTLC = []uint8 {0, 1, 3}
var chkTRC = []uint8 {1, 2, 4}
var chkBLC = []uint8 {3, 5, 6}
var chkBRC = []uint8 {4, 6, 7}

// This is where they live
// Must be initialized after specifying Width and Heigh
type World struct {
	// World size data
	Width uint
	Heigh uint
	Size uint

	// World corners
	botRC uint
	botLC uint
	topRC uint
	topLC uint

	// Where to move to check for neighbours
	cellNeighbours [8]int

	cells []bool
}


// ###############    PUBLIC METHODS    ###############

// Creates the slice of cells based on worlds width and heigh
// and gives values to important positions in the world
func (w *World) Initialize() {
	w.Size = w.Width * w.Heigh

	w.cells = make([]bool, w.Size)

	w.topLC = 0
	w.topRC = w.Width - 1
	w.botLC = (w.Size - 1) - (w.Width - 1)
	w.botRC = w.Size - 1

	w.cellNeighbours = [8]int {
		int(-w.Width -1), int(-w.Width), int(-w.Width + 1),
		-1,		 			    		   				1 ,
		int(w.Width - 1), int(w.Width),  int(w.Width +  1),
	}
}

// Returns cell state
func (w *World) GetCell(x, y uint) (bool, error) {
	pos := w.toAbs(x, y)

	if err := w.checkCoords(pos); err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return w.cells[pos], nil
}

// Sets cell state
func (w *World) SetCell(x, y uint, state bool) error {
	pos := w.toAbs(x, y)

	if err := w.checkCoords(pos); err != nil {
		return fmt.Errorf("%w", err)
	}

	w.cells[pos] = state

	return nil
}

// Calculate alive neighbours
func (w *World) AliveNeighbours(x, y uint) (uint8, error) {
	pos := w.toAbs(x, y)

	err := w.checkCoords(pos)
	if err != nil {
		return 255, err
	}

	acc := 0
	switch w.checkCase(pos) {
	case lefE:
		return w.countNeighbours(chkLE, pos), nil
	case rigE:
		return w.countNeighbours(chkRE, pos), nil
	case topE:
		return w.countNeighbours(chkTE, pos), nil
	case botE:
		return w.countNeighbours(chkBE, pos), nil
	case topLC:
		return w.countNeighbours(chkTLC, pos), nil
	case botLC:
		return w.countNeighbours(chkBLC, pos), nil
	case topRC:
		return w.countNeighbours(chkTRC, pos), nil
	case botRC:
		return w.countNeighbours(chkBRC, pos), nil
	case inCell:
		return w.countNeighbours(chkF, pos), nil
	}

	return uint8(acc), nil
}


// ###############    PRIVATE METHODS    ###############

// Converts (x, y) coordinates to absolute position 
func (w *World) toAbs(x, y uint) uint {
	return y * w.Width + x
}

// Checks if coordinates are inside the world
func (w *World) checkCoords(pos uint) error {
	if pos > w.Size {
		return fmt.Errorf("Cell coordinates out of bounds for a %dx%d world", w.Width, w.Heigh)
	}

	return nil
}

// Returns relative position of cell like top left corner
// Does not check for error as its called always after checking
func (w *World) checkCase(pos uint) int8 {
	switch pos {
	case w.topLC:
		return topLC
	case w.topRC:
		return  topRC
	case w.botLC:
		return botLC
	case w.botRC:
		return botRC
	default:
		if pos < w.Width{
			return topE
		} else if pos > w.botLC && pos < w.Size {
			return botE									
		} else {
			if rest := pos % w.Width; rest == 0 {
				return lefE
			} else if rest == w.Width - 1 {
				return rigE
			} else {
				return inCell
			}
		}

	}
}

// Returns the number of the check type neighbours
// Does not check for error as its called always after checking
func (w *World) countNeighbours(check []uint8, pos uint) uint8 {
	acc := 0

	for _, v := range check {
		if w.cells[uint(int(pos) + w.cellNeighbours[v])] {
			acc += 1
		}
	}

	return uint8(acc)
}
