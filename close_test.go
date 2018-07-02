package wporg

import (
	"os"
	"testing"
)

// These tests need revisiting. They need to cover a range of potential events
// to properly test the Drain/Close checks

func TestCheckClose(t *testing.T) {
	var err error

	f, err := os.Open("testdata/close.txt")
	if err != nil {
		t.Error("Could not open testdata/close.txt")
	}

	checkClose(f, &err)

	got := err
	expect := "" // nil value

	if nil != got {
		t.Errorf("Expected %#v got %#v", expect, got)
	}
}

func TestCheckDrainAndClose(t *testing.T) {
	var err error

	f, err := os.Open("testdata/close.txt")
	if err != nil {
		t.Error("Could not open testdata/close.txt")
	}

	drainAndClose(f, &err)

	got := err
	expect := "" // nil value

	if nil != got {
		t.Errorf("Expected %#v got %#v", expect, got)
	}
}
