package sago

import (
	"math"
)

/*
simplex Algorithm Go (SAGo) is a simple utility that runs the simplex algorithm
on a given objective function in standard equality form (SEF)
See sample usage in Simplex_test > Test_USE_EXAMPLE
*/

//Feasible returns true if the LP is feasible, false otherwise
func (lp *LP) Feasible() bool {
	if !lp.feasibilityKnown {
		lp.simplexPhaseI()
	}
	return lp.feasible
}

//Optimal returns true if the LP solved and optimal
func (lp *LP) Optimal() bool {
	if lp.optimal {
		return lp.optimal
	}
	for _, i := range lp.tableau[0][1:] {
		if Flt(i, 0, EPSILON) {
			return false
		}
	}
	lp.optimal = true
	return true
}

//Optimize executes the simplex algorithm on the LP
func (lp *LP) Optimize() {
	if lp.objective == MAXIMIZE {
		lp.negateObjectiveFunction()
	}
	lp.simplexPhaseI()
	if lp.objective == MINIMIZE {
		lp.negateObjectiveFunction()
	}
	if !lp.feasible {
		lp.unbounded = false
		lp.optimal = false
		lp.solved = true
		return
	}
	lp.simplex()
}

func (lp *LP) simplexPhaseI() {
	A := lp.auxLP()
	A.auxPreSimplex()
	A.simplex()
	if !lp.feasibilityKnown {
		lp.feasibilityKnown = true
		lp.feasible = Feq(A.ObjectiveValue(), 0, EPSILON)
	}
	l := len(lp.tableau[0])
	for i := range lp.tableau[1:] {
		lp.tableau[i+1] = A.tableau[i+1][:l]
	}
}

func (lp *LP) auxPreSimplex() {
	n := lp.NumConstraints()
	for r := 1; r <= n; r++ {
		for c := range (*lp).tableau[0] {
			(*lp).tableau[0][c] -= (*lp).tableau[r][c]
		}
	}
}

func (lp *LP) simplex() {
	for !lp.solved {
		lp.simplexIteration()
	}
}

func (lp *LP) simplexIteration() {
	if lp.solved {
		return
	}

	// Test optimal
	if lp.Optimal() {
		lp.solved = true
		lp.unbounded = false
		return
	}

	// Choose column to perform ratio test
	col := 1
	for i, j := range lp.tableau[0][1:] {
		if Flt(j, 0, EPSILON) {
			col = i + 1
			break
		}
	}

	// Ratio test
	row := lp.ratioTest(col)

	// unbounded
	if row == 0 {
		lp.unbounded = true
		lp.solved = true
		lp.optimal = false
		return
	}

	// Set the selected variable to 1 by multiplying row by the inverse of variable
	inverse := 1 / (lp.tableau)[row][col]
	for i := range (lp.tableau)[row] {
		(lp.tableau)[row][i] *= inverse
	}

	// make all other values in column 0
	for i, r := range lp.tableau {
		factor := r[col]
		if Feq(factor, 0, EPSILON) || i == row {
			continue
		}
		for j := range r {
			(lp.tableau)[i][j] -= factor * (lp.tableau)[row][j]
		}
	}
}

// returns the index of the constraint with the lowest ratio
// index 1 is the first constraint
func (lp *LP) ratioTest(column int) int {
	row := -1
	low := math.MaxFloat64
	for i, c := range lp.ListConstraints() {
		if Fle(c[column], 0, EPSILON) { // divide by 0
			continue
		}
		ratio := c[0] / c[column]
		if Fle(0, ratio, EPSILON) && Flt(ratio, low, EPSILON) {
			row, low = i, ratio
		}
	}
	return row + 1
}
