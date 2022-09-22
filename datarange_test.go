// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package datarange

import "testing"

func TestNew(t *testing.T) {

	dr1 := Make(10.123456, 100.98765, 1, "meter")
	if dr1.low != 10.1 || dr1.high != 101.0 {
		t.Errorf("RoundingRange fail: want low 10.1 high 101, get %v, %v", dr1.low, dr1.high)
	}

	dr2 := Make(10.123456, 100.98765, -1, "meter")
	if dr2.low != 10 || dr2.high != 110 {
		t.Errorf("RoundingRange fail: want low 10 high 110, get %v, %v", dr2.low, dr2.high)
	}

	dr3 := Make(10.123456, 100.98765, 3, "meter")
	if dr3.low != 10.123 || dr3.high != 100.988 {
		t.Errorf("RoundingRange fail: want low 10.123 high 100.988, get %v, %v", dr3.low, dr3.high)
	}

	dr4 := Make(10.123456, 100.98765, 0, "meter")
	if dr4.low != 10 || dr4.high != 101 {
		t.Errorf("RoundingRange fail: want low 10 high 101, get %v, %v", dr4.low, dr4.high)
	}
}

func TestProgress(t *testing.T) {

	dr1 := Make(0, 100, 1, "meter")
	r := dr1.Progress(10)
	if r != 0.1 {
		t.Errorf("Progress fail: want 0.1, get %v", r)
	}

	r = dr1.Progress(-1)
	if r != 0 {
		t.Errorf("Progress fail: want 0, get %v", r)
	}

	r = dr1.Progress(110)
	if r != 1 {
		t.Errorf("Progress fail: want 1, get %v", r)
	}

	dr2 := Make(10000, 20000, 1, "meter")
	r = dr2.Progress(12000)
	if r != 0.2 {
		t.Errorf("Progress fail: want 0.2, get %v", r)
	}
}

func FuzzBuildAutoStepsize(f *testing.F) {

	f.Add(0.0, 10.0, 10.0)
	f.Add(3.0, 7.0, 10.0)
	f.Add(5.0, 5.0, 10.0)
	f.Add(8.0, 2.0, 10.0)
	f.Add(8.0, -2.0, 10.0)
	f.Add(9.0, 1.0, 10.0)
	f.Add(1000.0, 64000.0, 10.0)
	f.Add(12456.0, 45789.0, 20.0)
	f.Add(1.0, 10.0, 3.0)
	f.Add(-10.0, 10.0, 100.0)
	f.Add(0.002587, 1.0, 12.0)
	f.Add(19421.8139685769, 980.8445405, 10.0)
	f.Add(109421.8139, 800018.8445405, 5000.0)

	f.Fuzz(func(t *testing.T, l float64, d float64, ms float64) {
		dr := Make(l, l+d, -ms, "fuzz")
		if dr.Steps() == 0 || dr.Steps() > uint(-ms) {
			t.Errorf("l:%v h:%v maxsteps:%v --> %v nbSteps:%v", l, l+d, ms, dr, dr.Steps())
		}

	})

}
