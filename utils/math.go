package utils

import (
	"math"
)

const float64EqualityThreshold = 1e-9

func MinInt(vn ...int) int {
	m := vn[0]
	for _, val := range vn {
		if val < m {
			m = val
		}
	}
	return m
}

// MinIntSlice will panic if slice is empty.
func MinIntSlice(ss []int) int {
	return MinInt(ss...)
}

func MinInt64(vn ...int64) int64 {
	m := vn[0]
	for _, val := range vn {
		if val < m {
			m = val
		}
	}
	return m
}

// MinInt64Slice will panic if slice is empty.
func MinInt64Slice(ss []int64) int64 {
	return MinInt64(ss...)
}

func MaxInt(vn ...int) int {
	m := vn[0]
	for _, val := range vn {
		if val > m {
			m = val
		}
	}
	return m
}

func MaxIntSlice(ss []int) int {
	return MaxInt(ss...)
}

func MaxInt64(vn ...int64) int64 {
	m := vn[0]
	for _, val := range vn {
		if val > m {
			m = val
		}
	}
	return m
}

func MaxInt64Slice(ss []int64) int64 {
	return MaxInt64(ss...)
}

// SumInt also return int64 because of the sum may over int.
func SumInt(vn ...int) int64 {
	sum := int64(0)

	for _, val := range vn {
		sum += int64(val)
	}

	return sum
}

func SumIntSlice(ss []int) int64 {
	return SumInt(ss...)
}

func SumInt64(vn ...int64) int64 {
	sum := int64(0)

	for _, val := range vn {
		sum += val
	}

	return sum
}

func SumInt64Slice(ss []int64) int64 {
	return SumInt64(ss...)
}

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
