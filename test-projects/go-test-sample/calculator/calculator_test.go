package calculator

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}

func TestAddNegative(t *testing.T) {
	result := Add(-1, -1)
	if result != -2 {
		t.Errorf("Add(-1, -1) = %d; want -2", result)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	if result != 2 {
		t.Errorf("Subtract(5, 3) = %d; want 2", result)
	}
}

func TestMultiply(t *testing.T) {
	result := Multiply(4, 3)
	if result != 12 {
		t.Errorf("Multiply(4, 3) = %d; want 12", result)
	}
}
