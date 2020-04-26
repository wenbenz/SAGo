package sago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLP_SetObjectiveFunctionMax(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 3, 1, 2)
	assert.Equal(t, lp.ObjectiveFunctionVector(), []float64{3, 1, 2})
	assert.Equal(t, lp.tableau[0], []float64{3, 1, 2})
}

func TestLP_AddConstraintEq(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintEq(3, 1, 2)
	assert.Equal(t, tableau{{1, 2, 3}, {3, 1, 2}}, lp.tableau)
	lp.AddConstraintEq(4, 5, 6)
	assert.Equal(t, tableau{{1, 2, 3}, {3, 1, 2}, {4, 5, 6}}, lp.tableau)

	lp = NewLP()
	lp.AddConstraintEq(1, 2, 3)
	assert.Equal(t, tableau{{0, 0, 0}, {1, 2, 3}}, lp.tableau)
	lp.AddConstraintEq(4, 5, 6)
	assert.Equal(t, tableau{{0, 0, 0}, {1, 2, 3}, {4, 5, 6}}, lp.tableau)

	lp = NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintEq(3, 1, 2, 4)
	assert.Equal(t, tableau{{1, 2, 3, 0}, {3, 1, 2, 4}}, lp.tableau)
	lp.AddConstraintEq(3, 1, 2, 4, 6, 7, 8)
	assert.Equal(t, tableau{{1, 2, 3, 0, 0, 0, 0}, {3, 1, 2, 4, 0, 0, 0}, {3, 1, 2, 4, 6, 7, 8}}, lp.tableau)

	lp = NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintEq(-3, 1, 2, 4)
	assert.Equal(t, tableau{{1, 2, 3, 0}, {3, -1, -2, -4}}, lp.tableau)
}

func TestLP_AddConstraintLeq(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintGeq(3, 1, 2)
	assert.Equal(t, tableau{{1, 2, 3, 0}, {3, 1, 2, 1}}, lp.tableau)
	lp.AddConstraintGeq(4, 5, 6)
	assert.Equal(t, tableau{{1, 2, 3, 0, 0}, {3, 1, 2, 1, 0}, {4, 5, 6, 0, 1}}, lp.tableau)

	lp = NewLP()
	lp.AddConstraintGeq(1, 2, 3, 1)
	assert.Equal(t, tableau{{0, 0, 0, 0, 0}, {1, 2, 3, 1, 1}}, lp.tableau)
	lp.AddConstraintGeq(4, 5, 6)
	assert.Equal(t, tableau{{0, 0, 0, 0, 0, 0}, {1, 2, 3, 1, 1, 0}, {4, 5, 6, 0, 0, 1}}, lp.tableau)
}

func TestLP_AddConstraintGeq(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintLeq(3, 1, 2)
	assert.Equal(t, tableau{{1, 2, 3, 0}, {3, 1, 2, -1}}, lp.tableau)
	lp.AddConstraintLeq(4, 5, 6)
	assert.Equal(t, tableau{{1, 2, 3, 0, 0}, {3, 1, 2, -1, 0}, {4, 5, 6, 0, -1}}, lp.tableau)

	lp = NewLP()
	lp.AddConstraintLeq(1, 2, 3, 1)
	assert.Equal(t, tableau{{0, 0, 0, 0, 0}, {1, 2, 3, 1, -1}}, lp.tableau)
	lp.AddConstraintLeq(4, 5, 6)
	assert.Equal(t, tableau{{0, 0, 0, 0, 0, 0}, {1, 2, 3, 1, -1, 0}, {4, 5, 6, 0, 0, -1}}, lp.tableau)
}

func TestLP_ListConstraints(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintEq(3, 1, 2)
	lp.AddConstraintEq(7, 8, 9)
	lp.AddConstraintEq(4, 5, 6)
	assert.Equal(t, [][]float64{{3, 1, 2}, {7, 8, 9}, {4, 5, 6}}, lp.ListConstraints())
}

func TestLP_RemoveConstraint(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintEq(3, 1, 2)
	lp.AddConstraintEq(7, 8, 9)
	lp.AddConstraintEq(4, 5, 6)
	lp.RemoveConstraint(1)
	assert.Equal(t, tableau{{1, 2, 3}, {3, 1, 2}, {4, 5, 6}}, lp.tableau)
	lp.RemoveConstraint(0)
	assert.Equal(t, tableau{{1, 2, 3}, {4, 5, 6}}, lp.tableau)
	lp.RemoveConstraint(0)
	assert.Equal(t, tableau{{1, 2, 3}}, lp.tableau)
}

func TestLP_ClearConstraints(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 1, 2, 3)
	lp.AddConstraintEq(3, 1, 2)
	lp.AddConstraintEq(7, 8, 9)
	lp.AddConstraintEq(4, 5, 6)
	lp.ClearConstraints()
	assert.Equal(t, tableau{{1, 2, 3}}, lp.tableau)
}

func TestLP_AuxLP(t *testing.T) {
	lp := readLP("LP_feas_sef")
	expected := tableau{
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1},
		{4, 1, 1, 2, 1, 0, 0, 1, 0, 0},
		{5, 2, 0, 3, 0, 1, 0, 0, 1, 0},
		{7, 2, 1, 3, 0, 0, 1, 0, 0, 1},
	}
	actual := lp.auxLP().tableau
	assert.Equal(t, expected, actual)
}

func TestLP_DualLP(t *testing.T) {
	lp := readLP("LP_feas_min")
	dlp := readLP("LP_feas_min_dual")
	dlpMatrix := tableau{
		{0, 11, 5, 0, 35},
		{4, 1, 1, -1, 12},
		{5, 1, -1, -1, 0},
		{6, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0},
	}
	actual := lp.DualLP()
	assert.Equal(t, dlpMatrix, actual.tableau)
	assert.Equal(t, dlp, actual)
}

func TestLP_Solution(t *testing.T) {
	lp := readLP("LP_solved")
	lp.solved = true
	lp.optimal = true
	lp.feasible = true
	actual, _ := lp.Solution()
	assert.Equal(t, []float64{2.5, 1.5, 0, 0, 0, .5}, actual)

	lp.solved = false
	actual, err := lp.Solution()
	assert.Equal(t, []float64{}, actual)
	assert.Error(t, err)

	lp.solved = true
	lp.optimal = false
	lp.unbounded = true
	actual, err = lp.Solution()
	assert.Equal(t, []float64{}, actual)
	assert.Error(t, err)

	lp = readLP("LP_solved_3x9")
	lp.solved = true
	lp.optimal = true
	lp.feasible = true
	actual, _ = lp.Solution()
	assert.Equal(t, []float64{65.76374332608265, 0, 0, 0, 10.038553835008601, 0, 0, 0}, actual)
}
