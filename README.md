# datarange package in go

DataRange represents a range bounded by a low and a high value. 

[![Go Reference](https://pkg.go.dev/badge/github.com/sunraylab/datarange.svg)](https://pkg.go.dev/github.com/sunraylab/datarange)

stepsize property allows stepping througth the datarange. If stepsize is defined then boundaries are rounded at stepsize level to ensure an integer number of steps.

stepsize can be automatically calculated according to a requested maximum number of steps. Calculated StepSize is a power of 1.0, 2.5 and 5.0. For example 100, 250, 500, 5000, 50000, or 0.25, 0.1 are calculated stepsize.

This is very usefull to build axis scale on a chart for example.

## Usage

A datarange must be created with a factory

```go
    // simple data range with values from 0 to 10 meters, without stepsize
	dr1 := Make(0, 10, 0, "meter")
    // simple data range with values from 0 to 10 meters, with 10 steps of 1 meter
	dr2 := Make(0, 10, 1, "meter")
    // data range with values from 1 to 10 meters, leaving the factory calculating the best stepsize to get a maximum number of 20 steps
	dr3 := Make(0, 10, -20, "meter")
```

see examples in the [DataRange package documentation](https://pkg.go.dev/github.com/sunraylab/datarange#pkg-examples)

## Installing

```bash
go get -u github.com/sunraylab/datarange@latest
```

## Changelog

- v1.2.0: migration to larry868 & go v1.23
- v1.1.0: fix some bugs and rename func ``Build`` to ``Make`` to follow go naming guidelines
- v1.0.0: first release

