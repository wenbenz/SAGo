package sago

import "fmt"

const (
	//MAXIMIZE objective for LPs
	MAXIMIZE = iota
	//MINIMIZE objective for LPs
	MINIMIZE
)

//LP is a data structure that represents a linear program
type LP struct {
	feasible         bool
	feasibilityKnown bool
	objective        int // MAXIMIZE or MINIMIZE
	optimal          bool
	solved           bool
	tableau          tableau
	unbounded        bool
	width            int
}

type tableau [][]float64

//NewLP reates a new LP
func NewLP() *LP {
	return &LP{
		tableau:          tableau{[]float64{}},
		optimal:          false,
		solved:           false,
		unbounded:        false,
		feasible:         false,
		feasibilityKnown: false,
		objective:        MAXIMIZE,
		width:            0,
	}
}

//GetObjective returns the objective of the LP (one of { MAXIMIZE, MINIMIZE })
func (lp *LP) GetObjective() int {
	return lp.objective
}

//Solution returns the solution vector of the LP i.e. the value X=(x1, x2, ...)
//for the objective function MAX/MIN c1 x1 + x2 x2 + ...
func (lp *LP) Solution() ([]float64, error) {
	if !lp.solved {
		return []float64{}, SolutionUnavailableError{"LP is unsolved; try optimizing LP first!"}
	}
	if lp.unbounded {
		return []float64{}, NoSolutionError{"LP is unbounded"}
	}
	if !lp.feasible {
		return []float64{}, NoSolutionError{"LP is infeasible"}
	}
	var sol []float64
	var vals []float64

	// b
	for _, row := range lp.tableau[1:] {
		vals = append(vals, row[0])
	}

	for c, v := range lp.tableau[0][1:] {
		if v != 0 {
			// Skip trying to find a 1 in column; this entry in the solution vector is 0
			goto AddZero
		}

		// Try to find 1 in column
		for r := range lp.tableau[1:] {
			if vals == nil {
				return []float64{}, fmt.Errorf("an unexpected error has occurred")
			}
			if Feq(lp.tableau[r+1][c+1], 1, EPSILON) {
				// found 1, add value to solution vector, then skip to next iteration of outer loop
				sol = append(sol, vals[r])
				goto Next
			}
		}

	AddZero:
		sol = append(sol, 0)

	Next:
	}

	return sol, nil
}

//ObjectiveFunctionVector returns the vector C=(c1, c2, ...) in the objective
//function MAX/MIN c1 x1 + x2 x2 + ...
func (lp *LP) ObjectiveFunctionVector() []float64 {
	return lp.tableau[0]
}

//ObjectiveValue is the objective value of the *current* LP
func (lp *LP) ObjectiveValue() float64 {
	return (lp.tableau)[0][0]
}

//SetObjectiveFunction sets the objective function
//MAXIMIZE/MINIMIZE z = c1x1 + c2x2 + ... for the LP
func (lp *LP) SetObjectiveFunction(objective int, z float64, coefficients ...float64) {
	lp.tableau[0] = append([]float64{z}, coefficients...)
	lp.objective = objective
	fnLen := len(lp.tableau[0])
	if fnLen > lp.width {
		lp.increaseWidth(fnLen)
	}
}

//NumConstraints returns the number of constraints in the LP
func (lp *LP) NumConstraints() int {
	return len(lp.tableau) - 1
}

//AddConstraintEq  adds an equality constraint
//for b = a1x1 + a2x2 + ... + anxn, coefficients should be (a1, a2, ..., an)
func (lp *LP) AddConstraintEq(b float64, coefficients ...float64) {
	constraint := extend(append([]float64{b}, coefficients...), lp.width)
	if Flt(b, 0, EPSILON) {
		constraint = ScalarVectorMultiply(-1, constraint)
	}
	if len(constraint) > lp.width {
		lp.increaseWidth(len(constraint))
	}
	lp.tableau = append(lp.tableau, constraint)
}

//AddConstraintGeq is like AddConstraintEq but for >= constraints
func (lp *LP) AddConstraintGeq(b float64, coefs ...float64) {
	constraint := extend(coefs, lp.width-1)
	constraint = append(constraint, 1)
	lp.AddConstraintEq(b, constraint...)
}

//AddConstraintLeq is like AddConstraintEq but for <= constraints
func (lp *LP) AddConstraintLeq(b float64, coefs ...float64) {
	constraint := extend(coefs, lp.width-1)
	constraint = append(constraint, -1)
	lp.AddConstraintEq(b, constraint...)
}

//ListConstraints returns the tail of the tableau
func (lp *LP) ListConstraints() [][]float64 {
	return (lp.tableau)[1:]
}

//RemoveConstraint removes the i-th constraint from the LP
func (lp *LP) RemoveConstraint(i int) {
	lp.tableau = append((lp.tableau)[:i+1], (lp.tableau)[i+2:]...)
}

//ClearConstraints removes all constraints
func (lp *LP) ClearConstraints() {
	lp.tableau = (lp.tableau)[0:1]
}

//DualLP returns a pointer to the dual of the LP
func (lp *LP) DualLP() *LP {
	dlp := NewLP()
	dlp.tableau = Transpose(lp.tableau)
	return dlp
}

//Solved returns true if the simplex algorithm has been run on the LP
func (lp *LP) Solved() bool {
	return lp.solved
}

//Bounded returns true if the LP is bounded; false otherwise
func (lp *LP) Bounded() bool {
	return !lp.unbounded
}

//Width returns the width of the tableau
func (lp *LP) Width() int {
	return len(lp.tableau[0])
}

func extend(A []float64, length int) []float64 {
	for len(A) < length {
		A = append(A, 0)
	}
	return A
}

func (lp *LP) negateObjectiveFunction() {
	lp.tableau[0] = ScalarVectorMultiply(-1, lp.tableau[0])
	lp.tableau[0][0] *= -1
}

func (lp *LP) increaseWidth(length int) {
	lp.width = length
	for i := range lp.tableau {
		lp.tableau[i] = extend(lp.tableau[i], length)
	}
}

func (lp *LP) auxLP() *LP {
	variables := len(lp.ObjectiveFunctionVector())
	constraints := lp.ListConstraints()
	var objFunc []float64

	// Auxiliary LP obj function
	for i := 0; i < variables; i++ {
		objFunc = append(objFunc, 0)
	}
	for range constraints {
		objFunc = append(objFunc, 1)
	}

	// Start building LP
	matrix := NewLP()
	matrix.tableau = tableau{objFunc}

	for i := range constraints {
		matrix.tableau = append(matrix.tableau, constraints[i])
		for range constraints {
			matrix.tableau[i+1] = append(matrix.tableau[i+1], 0)
		}
		matrix.tableau[i+1][variables+i] = 1
	}

	return matrix
}
