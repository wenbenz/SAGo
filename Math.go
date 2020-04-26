package sago

//EPSILON is 1/2^32
const EPSILON = float64(1) / float64(2<<32)

//Transpose returns the transpose of a matrix represented by a 2d array
func Transpose(matrix [][]float64) [][]float64 {
	yl := len(matrix)
	if yl == 0 {
		return [][]float64{}
	}
	xl := len(matrix[0])
	if xl == 0 {
		return [][]float64{{}}
	}
	result := make([][]float64, xl)
	for i := range result {
		result[i] = make([]float64, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = matrix[j][i]
		}
	}
	return result
}

//ScalarVectorMultiply multiplies a vector by a scalar factor
func ScalarVectorMultiply(factor float64, vector []float64) []float64 {
	A := vector
	for i := range A {
		A[i] *= factor
	}
	return A
}

//Flt float less than
func Flt(f1, f2, error float64) bool {
	return f1+error < f2
}

//Fle float less than or equal to
func Fle(f1, f2, error float64) bool {
	return f1-error < f2
}

//Fgt float greater than
func Fgt(f1, f2, error float64) bool {
	return f1-error > f2
}

//Fge float greater than or equal to
func Fge(f1, f2, error float64) bool {
	return f1+error > f2
}

//Feq Float equals
func Feq(f1, f2, error float64) bool {
	return Fle(f1, f2, error) && Fge(f1, f2, error)
}
