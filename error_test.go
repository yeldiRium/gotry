package gotry

import "testing"

func TestIsTimeout(t *testing.T) {
	if !IsTimeout(ErrTimeout) {
		t.Error("ErrTimeout is somehow broken.")
	}
}

func TestIsMaxTriesReached(t *testing.T) {
	if !IsMaxTriesReached(ErrMaxTriesReached) {
		t.Error("ErrMaxTriesReached is somehow broken.")
	}
}

func TestIsFIsMissing(t *testing.T) {
	if !IsFIsMissing(ErrFIsMissing) {
		t.Error("ErrFIsMissing is somehow broken.")
	}
}
