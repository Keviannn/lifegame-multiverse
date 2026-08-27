package rules

import "testing"

func TestApplyConway(t *testing.T) {
	var r Rules
	r.Parse("B3/S23")
	b := [9]bool {false, false, false, true, false, false, false, false, false}
	s := [9]bool {false, false, true, true, false, false, false, false, false}

	for i, v := range b {
		if r.birth[i] != v {
			t.Errorf("Birth is not parsed well")
		}
	}

	for i, v := range s {
		if r.survive[i] != v {
			t.Errorf("Survive is not parsed well")
		}
	}
}

func TestInvalidLetters(t *testing.T) {
	var r Rules
	err := r.Parse("xyz")
	if err == nil {
		t.Errorf("Expected error")
	}
	t.Logf("Intentional Error: %v", err)
}

func TestInvalidNumbers(t *testing.T) {
	var r Rules
	err := r.Parse("B3/S29")
	if err == nil {
		t.Errorf("Expected error")
	}
	t.Logf("Intentional Error: %v", err)
}

func TestApplyRules(t *testing.T) {
	var r Rules
	r.Parse("B3/S23")
	v := r.DecideFate(false, 3)
	if !v {
		t.Error("Cell expected to be born")
	}

	v = r.DecideFate(true, 1)
	if v {
		t.Error("Cell expected to die")
	}
}

