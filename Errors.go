package sago

import "fmt"

//InvalidInputError input is invalid
type InvalidInputError struct {
	s           string
}

func(e InvalidInputError)Error()string{return fmt.Sprintf("%s: %s", "Invalid input", e.s)
}

//NoSolutionError there is no solution to this LP
type NoSolutionError struct {
	s string
}

func (e NoSolutionError) Error() string {
	return fmt.Sprintf("%s: %s", "No solution", e.s)
}

//SolutionUnavailableError the solution has not yet been calculated
type SolutionUnavailableError struct {
	s string
}

func (e SolutionUnavailableError) Error() string {
	return fmt.Sprintf("%s: %s", "Solution unavailable", e.s)
}
