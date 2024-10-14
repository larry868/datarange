// Copyright 2022-2024 by larry868. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

/*
datarange represents a range bounded by two values low and high. Stepsize property allows Boundaries to be rounded and to calculate steps.
*/
package datarange

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// A DataRange bounded with two float64 values, rounded at stepsize level and expressed in a spcecific unit.
//
// Use the DataRange factory to create a new one
type DataRange struct {
	unit     string  // define the unit of low and high boundaries. For information
	stepsize float64 // zero means no step size
	low      float64 // low boundary
	high     float64 // high boundary
}

// default string formating: "unit[ low :stepsize: high ]"
func (thisdr DataRange) String() string {
	var strlow, strhigh string
	if thisdr.stepsize != 0 {
		strlow = FormatData(thisdr.low, thisdr.stepsize)
		strhigh = FormatData(thisdr.high, thisdr.stepsize)
	} else {
		strlow = fmt.Sprintf("%v", thisdr.low)
		strhigh = fmt.Sprintf("%v", thisdr.high)
	}
	strprec := fmt.Sprintf("%v", thisdr.stepsize)
	return fmt.Sprintf("%s[ %s :%s: %s ]", thisdr.unit, strlow, strprec, strhigh)
}

// Returns a new Datarange, with boudaries rounded at stepsize level:
//   - if stepsize > 1 then stepsize is the rounding number of digit.
//   - if stepsize < 1 then stepsize is the rounding number of decimals.
//   - if stepsize == 0 then do not round
//   - if stepsize <= 0 then its absolute value is considered as the maxsteps and the stepsize is calculated automatically
//
// # The calculated StepSize is a power of 1.0 2.5 and 5.0, for example 100 250 500 5000 50000 or 0.25 0.1
//
// lowboudaries and highboundaries are automatically determined according to a and b values
//
// lowboundary is rounded down, highboundary is rounded up
func Make(a float64, b float64, stepsize float64, unit string) DataRange {
	dr := &DataRange{unit: unit}

	if stepsize >= 0 {
		dr.stepsize = stepsize
	} else {
		maxsteps := -stepsize
		delta := b - a
		if delta < 0 {
			delta = -delta
		}
		rawstep := delta / float64(maxsteps)
		exp := float64(int(math.Log10(rawstep)+1) - 1)

		basicstepsizes := []float64{1.0, 2.5, 5.0, 10.0, 25.0, 50.0}
		for _, basicstepsize := range basicstepsizes {
			possiblestepsize := basicstepsize * math.Pow(10, exp)

			roundeda := math.Floor(a/possiblestepsize) * possiblestepsize
			roundedb := math.Ceil(b/possiblestepsize) * possiblestepsize
			roundeddelta := roundedb - roundeda

			if possiblestepsize*float64(maxsteps) >= roundeddelta {
				dr.stepsize = possiblestepsize
				break
			}
		}
	}

	dr.ResetBoundaries(a, b)
	return *dr
}

// ResetBoundaries sets new boundaries without changing the stepsize.
//
// lowBoundaries and highBoundaries are automatically determined according to a and b values
//
// lowBoundary is rounded down, highBoundary is rounded up
func (pdr *DataRange) ResetBoundaries(a float64, b float64) {
	if b > a {
		pdr.low = a
		pdr.high = b
	} else {
		pdr.low = b
		pdr.high = a
	}

	if pdr.stepsize != 0 {
		// round boundaries
		pdr.low = math.Floor(pdr.low/pdr.stepsize) * pdr.stepsize
		pdr.high = math.Ceil(pdr.high/pdr.stepsize) * pdr.stepsize

		// round decimals to avoid dust
		ratio := math.Pow(10, float64(decimals(pdr.stepsize)))
		pdr.low = math.Floor(pdr.low*ratio) / ratio
		pdr.high = math.Ceil(pdr.high*ratio) / ratio
	}
}

func (thisdr DataRange) Unit() string {
	return thisdr.unit
}

func (thisdr DataRange) Low() float64 {
	return thisdr.low
}

func (thisdr DataRange) High() float64 {
	return thisdr.high
}

func (thisdr DataRange) StepSize() float64 {
	return thisdr.stepsize
}

// Enlarge the datarange with boundaries extended by the coef.
// if coef is <1 then boudaries are reduced.
// if coef <= 0 then the boundaries are reset to 0
func (pdr *DataRange) Enlarge(coef float64) {
	var lb, hb float64
	if coef > 0 {
		lb = pdr.low / coef
		hb = pdr.high * coef
	}
	pdr.ResetBoundaries(lb, hb)
}

// Returns the number of steps between the boundaries according to the stepsize.
//
// returns +Inf if thists.stepsize == 0
func (dr DataRange) Steps() uint {
	f := (dr.high - dr.low) / dr.stepsize
	if f < 0 {
		f = -f
	}
	return uint(f)
}

// Delta returns the value betwwen boundaries.
// Returns 0 if HighBoundary == LowBoundary.
// Rteurns <0 if HighBoundary is lower than LowBoundary.
func (dr DataRange) Delta() float64 {
	return dr.high - dr.low
}

// Progress returns the rate of val within the datarange, compared with the LowBoundary
// Return 0 if val is equal or under the lowboundary
// Return 1 if val is equal or greater the highboundary
// Return 0.5 if the datarange is composed of a single value and val is that value
func (dr DataRange) Progress(val float64) float64 {
	rate := 0.5
	rng := dr.high - dr.low
	if rng > 0 {
		rate = (val - dr.low) / rng
		if rate < 0 {
			rate = 0
		} else if rate > 1 {
			rate = 1
		}
	} else if val < dr.low {
		rate = 0
	} else if val > dr.high {
		rate = 1
	}
	return rate
}

// Equal checks if 2 ranges have the sames boundaries, the same stepsize, and the same unit
func (thisdr DataRange) Equal(another DataRange) bool {
	return thisdr.low == another.low && thisdr.high == another.high && thisdr.stepsize == another.stepsize
}

// FormatData f according to a stepsize
//
//	If stepsize == 0 then f is formatted without trailing zeros
//
// # Example
//
//	FormatData(11.2, 0.1) == "11.2"
//	FormatData(11.2, 1) == "11"
//	FormatData(11.2, 5) == "10"
func FormatData(f float64, stepsize float64) (str string) {
	if stepsize != 0 {
		str = fmt.Sprintf("%.*f", decimals(stepsize), f)
	} else {
		str = fmt.Sprintf("%v", f)
	}
	return str
}

/*
 * utility
 */

// returns the number of significant decimals for f
func decimals(f float64) (d uint) {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		x := len(s) - i - 1
		if x >= 0 {
			d = uint(x)
		}
	}
	return d
}
