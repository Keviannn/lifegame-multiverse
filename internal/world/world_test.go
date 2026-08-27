package world

import "testing"

func TestNewWorld(t *testing.T) {
	w := NewWorld(3, 3)

	if w.Size != 9 {
		t.Errorf("Expected world size 9: %d", w.Size)
	}
	if w.topLC != 0 {
		t.Errorf("Expected top left corner 0: %d", w.topLC)
	}
	if w.topRC != 2 {
		t.Errorf("Expected top right corner 2: %d", w.topRC)
	}
	if w.botLC != 6 {
		t.Errorf("Expected bot left corner 6: %d", w.botLC)
	}
	if w.botRC != 8 {
		t.Errorf("Expected bot right corner 8: %d", w.botRC)
	}
}

func TestSetCell(t *testing.T) {
	w := NewWorld(3, 3)

	err := w.SetCell(4, 4, Alive)
	if err == nil {
		t.Error("Expected error out of bounds")
	}
	t.Logf("Intentional Error: %v", err)

	w.SetCell(2, 2, Alive)
	if b, _ := w.GetCell(2, 2); !b {
		t.Error("Expected cell to be alive")
	}

	w.SetCell(2, 2, Dead)
	if b, _ := w.GetCell(2, 2); b {
		t.Error("Expected cell to be dead")
	}
}

func TestGetCell(t *testing.T) {
	w := NewWorld(3, 3)

	w.cells[2 * w.Width + 1] = Alive

	if b, _ := w.GetCell(1, 2); !b {
		t.Error("Expected cell to be alive")
	}
}

func TestCheckCase(t *testing.T) {
	w := NewWorld(3, 3)

	if w.checkCase(w.toAbs(0, 0)) != topLC {
		t.Error("Expected top left corner")
	}
	if w.checkCase(w.toAbs(2, 0)) != topRC {
		t.Error("Expected top right corner")
	}
	if w.checkCase(w.toAbs(0, 2)) != botLC {
		t.Error("Expected bot left corner")
	}
	if w.checkCase(w.toAbs(2, 2)) != botRC {
		t.Error("Expected bot right corner")
	}

	if w.checkCase(w.toAbs(1, 0)) != topE {
		t.Error("Expected top edge")
	}
	if w.checkCase(w.toAbs(0, 1)) != lefE {
		t.Error("Expected left edge")
	}
	if w.checkCase(w.toAbs(2, 1)) != rigE {
		t.Error("Expected right edge")
	}
	if w.checkCase(w.toAbs(1, 2)) != botE {
		t.Error("Expected bot edge")
	}

	if w.checkCase(w.toAbs(1, 1)) != inCell {
		t.Error("Expected inner cell")
	}
}

func TestCountNeighbours(t *testing.T) {
	w := NewWorld(3, 3)

	w.SetCell(0, 0, Alive)
	w.SetCell(2, 2, Alive)
	if w.countNeighbours(chkF, w.toAbs(1, 1)) != 2 {
		t.Error("Expected 2 neighbours")
	}
}

func AliveNeighbours(t *testing.T) {
	w := NewWorld(3, 3)

	w.SetCell(0, 0, Alive)
	w.SetCell(2, 2, Alive)

	if v, _ := w.AliveNeighbours(1, 1); v != 2 {
		t.Error("Expected 2 neighbours")
	}

	if _, err := w.AliveNeighbours(25, 1); err == nil {
		t.Error("Expected an error")
	}
}
