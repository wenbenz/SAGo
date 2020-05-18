package sago

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func readLP(testName string) *LP {
	testDir := "Tests"
	testBytes, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", testDir, testName))
	testRows := strings.Split(strings.TrimSpace(string(testBytes)), "\n")
	test := NewLP()
	test.tableau = tableau{}
	for _, row := range testRows {
		r := strings.Split(row, ",")
		test.tableau = append(test.tableau, []float64{})
		for _, v := range r {
			f, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
			test.tableau[len(test.tableau)-1] = append(test.tableau[len(test.tableau)-1], f)
		}
	}
	return test
}

func TestLP_SimplexIteration(t *testing.T) {
	lp := readLP("LP_feas_sef")
	lp.negateObjectiveFunction()
	expected := tableau{
		{7.5, 0, -2, .5, 0, 1.5, 0},
		{1.5, 0, 1, .5, 1, -.5, 0},
		{2.5, 1, 0, 1.5, 0, .5, 0},
		{2, 0, 1, 0, 0, -1, 1},
	}
	lp.simplexIteration()
	assert.Equal(t, expected, lp.tableau)
	assert.False(t, lp.unbounded)
	assert.False(t, lp.solved)

	lp = readLP("LP_unbounded_trivial")
	lp.simplexIteration()
	assert.True(t, lp.unbounded)
	assert.True(t, lp.solved)
}

func TestLP_Feasible(t *testing.T) {
	lp := readLP("LP_feas_sef")
	assert.True(t, lp.Feasible())

	lp = readLP("LP_unbounded_trivial")
	assert.True(t, lp.Feasible())
}

func TestLP_Feasible2(t *testing.T) {
	lp := readLP("LP_infeasibility_rounding_error")
	assert.True(t, lp.Feasible())
}

func TestLP_Simplex(t *testing.T) {
	lp := readLP("LP_feas_sef")
	lp.negateObjectiveFunction()
	lp.simplex()
	expected := tableau{
		{10.5, 0, 0, 1.5, 2, .5, 0},
		{1.5, 0, 1, .5, 1, -.5, 0},
		{2.5, 1, 0, 1.5, 0, .5, 0},
		{0.5, 0, 0, -.5, -1, -.5, 1},
	}
	assert.True(t, lp.solved)
	assert.True(t, lp.optimal)
	assert.False(t, lp.unbounded)
	assert.True(t, lp.Feasible())
	assert.Equal(t, expected, lp.tableau)
	assert.Equal(t, 10.5, lp.ObjectiveValue())
}

func TestLP_Optimize(t *testing.T) {
	// Has optimal solution
	lp := readLP("LP_feas_sef")
	lp.Optimize()
	expected := tableau{
		{10.5, 0, 0, 1.5, 2, .5, 0},
		{1.5, 0, 1, .5, 1, -.5, 0},
		{2.5, 1, 0, 1.5, 0, .5, 0},
		{0.5, 0, 0, -.5, -1, -.5, 1},
	}
	assert.True(t, lp.solved)
	assert.True(t, lp.optimal)
	assert.False(t, lp.unbounded)
	assert.True(t, lp.Feasible())
	assert.Equal(t, expected, lp.tableau)
	assert.Equal(t, 10.5, lp.ObjectiveValue())

	// unbounded
	lp = readLP("LP_unbounded_trivial")
	lp.Optimize()
	assert.True(t, lp.solved)
	assert.False(t, lp.optimal)
	assert.True(t, lp.unbounded)
	assert.True(t, lp.Feasible())

	// Infeasible
	lp = readLP("LP_infeas")
	lp.Optimize()
	assert.True(t, lp.solved)
	assert.False(t, lp.optimal)
	assert.False(t, lp.unbounded)
	assert.False(t, lp.Feasible())
}

func TestLP_Optimize2(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 0, 3, 2, 4, 0, 0, 0) // maximize   0 = (3,2,4,0,0,0)x
	lp.AddConstraintEq(4, 1, 1, 2, 1, 0, 0)                // constraint 4 = (1,1,2,1,0,0)x
	lp.AddConstraintEq(5, 2, 0, 3, 0, 1, 0)
	lp.AddConstraintEq(7, 2, 1, 3, 0, 0, 1)

	lp.Optimize()
	expected := tableau{
		{10.5, 0, 0, 1.5, 2, .5, 0},
		{1.5, 0, 1, .5, 1, -.5, 0},
		{2.5, 1, 0, 1.5, 0, .5, 0},
		{0.5, 0, 0, -.5, -1, -.5, 1},
	}
	assert.True(t, lp.solved)
	assert.True(t, lp.optimal)
	assert.False(t, lp.unbounded)
	assert.True(t, lp.Feasible())
	assert.Equal(t, expected, lp.tableau)
	assert.Equal(t, 10.5, lp.ObjectiveValue())
	solution, _ := lp.Solution()
	assert.Equal(t, []float64{2.5, 1.5, 0, 0, 0, .5}, solution)
}

func TestLP_Optimize3(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MINIMIZE, 0, 3, 2, 0, 0)
	lp.AddConstraintEq(6, 2, 1, 1, 0)
	lp.AddConstraintEq(4, 1, 1, 0, 1)

	lp.Optimize()
	assert.True(t, lp.solved)
	assert.True(t, lp.optimal)
	assert.False(t, lp.unbounded)
	assert.True(t, lp.Feasible())
	assert.Equal(t, float64(10), lp.ObjectiveValue())
	solution, _ := lp.Solution()
	assert.Equal(t, []float64{2, 2, 0, 0}, solution)
}

func TestLP_Optimize4(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MINIMIZE, 0, 6, 3, 0, 0)
	lp.AddConstraintEq(1, 1, 1, 1, 0)
	lp.AddConstraintEq(1, 2, -1, 0, 1)

	lp.Optimize()
	assert.True(t, lp.solved)
	assert.True(t, lp.optimal)
	assert.False(t, lp.unbounded)
	assert.True(t, lp.Feasible())
	assert.Equal(t, float64(5), lp.ObjectiveValue())
	solution, _ := lp.Solution()
	x1 := float64(2) / 3
	x2 := float64(1) / 3
	assert.Equal(t, []float64{x1, x2, 0, 0}, solution)
}

func TestLP_Optimize5(t *testing.T) {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 100, 1, 1)
	lp.AddConstraintEq(100, 1, 1)
	lp.AddConstraintEq(50, 1, 0)
	lp.AddConstraintEq(40, 1, 0)
	lp.Optimize()
	assert.False(t, lp.Feasible())

	lp = NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 100, 1, 1)
	lp.AddConstraintEq(100, 1, 1)
	lp.AddConstraintEq(50, 1, 0)
	lp.AddConstraintEq(60, 1, 0)
	lp.Optimize()
	assert.False(t, lp.Feasible())

	lp = NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 100, 1, 1)
	lp.AddConstraintEq(100, 1, 1)
	lp.AddConstraintEq(50, 1, 0)
	lp.AddConstraintLeq(40, 1, 0)
	lp.Optimize()
	assert.True(t, lp.Feasible())

	lp = NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 100, 1, 1)
	lp.AddConstraintEq(100, 1, 1)
	lp.AddConstraintEq(50, 1, 0)
	lp.AddConstraintGeq(60, 1, 0)
	lp.Optimize()
	assert.True(t, lp.Feasible())
}

func TestLP_Optimize6(t *testing.T) {
	lp := readLP("LP_sol_0")
	lp.Optimize()
	sol, _ := lp.Solution()
	assert.Equal(t, []float64{0, 0, 7095.5198887422475, 0, 10000, 0}, sol)
}

// Test degen iteration
func TestLP_Optimize7(t *testing.T) {
	lp := readLP("LP_degenerate_iteration")
	lp.Optimize()
	sol, _ := lp.Solution()
	assert.Equal(t, []float64{0, 0, 0, 2}, sol)
}

/*
Example of how to optimize using simplex to solve:
Max 	    3 x_1 + 5 x_2
s.t.
	50  =           5 x_2
	100 >= 15 x_1 + 2 x_2
*/
func ExampleLP_Optimize() {
	lp := NewLP()
	lp.SetObjectiveFunction(MAXIMIZE, 0, 3, 5)
	lp.AddConstraintEq(50, 0, 5)
	lp.AddConstraintGeq(100, 9, 2)

	// Optimize
	lp.Optimize()
	fmt.Println("feasible:", lp.Feasible())
	fmt.Println("optimal:", lp.Optimal())
	fmt.Println("bounded:", lp.Bounded())
	fmt.Println("objectiveValue:", lp.ObjectiveValue())
	solution, _ := lp.Solution()
	fmt.Printf("solution: x1 = %.2f, x2 = %.2f", solution[0], solution[1])
}
