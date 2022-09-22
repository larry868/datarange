// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

package datarange

import (
	"fmt"
	"math"
)

func ExampleDataRange() {
	// this example illustrate how to build an automatic scale appart from a dataset

	// sample dataset
	datasetMinValue := 25.56
	datasetMaxValue := 10850.12

	// build a range, with an automatic stepsize to get a max of 10 steps on the scale
	dr := Make(datasetMinValue, datasetMaxValue, -10, "amount")
	fmt.Printf("dataset = %v\n", dr)

	// step througth the datarange
	fmt.Println("scale")
	for y := dr.Low(); y <= dr.High(); y += dr.StepSize() {
		fmt.Printf(" %s\n", FormatData(y, dr.StepSize()))
	}

	// Output:
	// dataset = amount[ 0 :2500: 12500 ]
	// scale
	//  0
	//  2500
	//  5000
	//  7500
	//  10000
	//  12500
}

func ExampleDataRange_Steps() {

	var drs [8]DataRange
	drs[0] = Make(1, 10, 1, "meter")
	drs[1] = Make(-10, 10, 0.1, "meter")
	drs[2] = Make(1.245, 2.4, 0, "meter")
	drs[3] = Make(12448, 548983, 25, "meter")
	drs[4] = Make(1, 2, 10, "meter")
	drs[5] = Make(1.5, 1.5, 1, "meter")
	drs[6] = Make(10, 10, 10, "meter")
	drs[7] = Make(19421.8139685769, 20402.658509423058, -10.0, "test")

	var stri string
	for _, dr := range drs {
		ui := dr.Steps()
		if ui == uint(math.Inf(1)) {
			stri = "infinite"
		} else {
			stri = fmt.Sprintf("%d", ui)
		}
		fmt.Printf("%s, intervals=%s\n", dr, stri)
	}

	// Output:
	// meter[ 1 :1: 10 ], intervals=9
	// meter[ -10.0 :0.1: 10.0 ], intervals=200
	// meter[ 1.245 :0: 2.4 ], intervals=infinite
	// meter[ 12425 :25: 549000 ], intervals=21463
	// meter[ 0 :10: 10 ], intervals=1
	// meter[ 1 :1: 2 ], intervals=1
	// meter[ 10 :10: 10 ], intervals=0
	// meter[ 19250 :250: 20500 ], intervals=5
}

func ExampleMake() {

	type sample struct {
		low      float64
		high     float64
		maxsteps float64
	}
	samples := []sample{{0, 10, 10}, {3, 10, 10}, {5, 10, 10}, {8, 10, 10}, {9, 10, 10}, {1000, 65000, 10}, {12456, 45789, 20}, {9925, 10401, 10}}

	for _, s := range samples {
		dr := Make(s.low, s.high, -s.maxsteps, "amount")
		fmt.Printf("range[%v %v]/%v --> %v with %v intervals\n", s.low, s.high, s.maxsteps, dr.String(), dr.Steps())
	}

	// Output:
	// range[0 10]/10 --> amount[ 0 :1: 10 ] with 10 intervals
	// range[3 10]/10 --> amount[ 3 :1: 10 ] with 7 intervals
	// range[5 10]/10 --> amount[ 5.0 :0.5: 10.0 ] with 10 intervals
	// range[8 10]/10 --> amount[ 8.00 :0.25: 10.00 ] with 8 intervals
	// range[9 10]/10 --> amount[ 9.0 :0.1: 10.0 ] with 10 intervals
	// range[1000 65000]/10 --> amount[ 0 :10000: 70000 ] with 7 intervals
	// range[12456 45789]/20 --> amount[ 10000 :2500: 47500 ] with 15 intervals
	// range[9925 10401]/10 --> amount[ 9900 :100: 10500 ] with 6 intervals
}
