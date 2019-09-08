package proxyblock

import (
	"testing"
	"time"
)

func TestBlock(t *testing.T) {
	cd := NewCountData(5)
	var testVal int
	// now we just need to push in 6 hits
	for j := 0; j <= 6; j++ {
		if cd.ShouldBlock() {
			testVal = j
			break
		}
	}
	if testVal != 6 {
		t.Error("Block did not happen as expected")
	}
}

// TestNextMinute puts some hits into the window of bucket1 and then sleeps until the one minute is over
// to ensure that we behave properly on minute 2
func TestNextMinute(t *testing.T) {
	cd := NewCountData(5)
	// fill up to the allowed max
	for j := 0; j <= 5; j++ {
		cd.ShouldBlock()
	}
	time.Sleep(62 * time.Second)
	if cd.ShouldBlock() {
		t.Error("NextMinute -- expected not to be blocked after 60 seconds")
	}
}

// TestStraddleMinute tests the case of us crossing the threshold at near the one minute reset or at some point where we stagger
// our hits across minutes.  This is a bit sketchy because we depend on the clock, but let's have a go.
func TestStraddleMinute(t *testing.T) {
	cd := NewCountData(5)
	// fill up to the allowed max
	for j := 0; j <= 3; j++ {
		cd.ShouldBlock()
	}
	time.Sleep(55)
	for j := 0; j <= 3; j++ {
		if cd.ShouldBlock() {
			t.Log("straddle", j)
			if j != 2 {
				t.Error("StraddleMinute - did not expected to be blocked yet")
			}
			break
		}
		time.Sleep(2)
	}
}
