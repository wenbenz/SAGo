package sago

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTanspose1(t *testing.T) {
	matrix := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}
	expected := [][]float64{
		{1, 3, 5},
		{2, 4, 6},
	}
	assert.Equal(t, expected, Transpose(matrix))
}

func TestTanspose2(t *testing.T) {
	matrix := [][]float64{}
	expected := [][]float64{}
	assert.Equal(t, expected, Transpose(matrix))
}

func TestTanspose3(t *testing.T) {
	matrix := [][]float64{{}}
	expected := [][]float64{{}}
	assert.Equal(t, expected, Transpose(matrix))
}

func TestTranspose4(t *testing.T) {
	matrix := [][]float64{
		{0, 4, 5, 6, 0, 0, 0},
		{11, 1, 1, 0, 1, 0, 0},
		{5, 1, -1, 0, 0, 1, 0},
		{0, -1, -1, 1, 0, 0, 0},
		{35, 12, 0, 0, 0, 1, 0},
	}
	expected := [][]float64{
		{0, 11, 5, 0, 35},
		{4, 1, 1, -1, 12},
		{5, 1, -1, -1, 0},
		{6, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0},
	}
	assert.Equal(t, expected, Transpose(matrix))
}
