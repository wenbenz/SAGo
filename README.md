# SAGo - Simplex Algorithm Go
[![CircleCi](https://circleci.com/gh/wenbenz/SAGo.svg?style=shield)](https://circleci.com/gh/wenbenz/SAGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/wenbenz/SAGo)](https://goreportcard.com/report/github.com/wenbenz/SAGo)

Simple implementation of the 2-phase Simplex Algorithm in Go.

## Add to your project
```go get github.com/wenbenz/SAGo```

## Usage
See use example in `SimplexTest.go > ExampleLP_Optimize`
- Construct a new LP by calling `lp := NewLP()`
- Call `lp.SetObjectiveFunction(z, a1, a2, ...)` to set the objective function to maximize where z = (a1, a2, ...)x
- Call `lp.AddConstraint[Eq/Leq/Geq](bi, a1, a2, ...)` to add a constraint where bi = (a1, a2, ...)x
- Call `Optimize` when you're ready to go!
