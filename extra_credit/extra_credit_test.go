package main

import (
	"testing"
)

func TestSwapIntegers(t *testing.T) {
	a, b := 1, 2
	swapTwoIntegers(&a, &b)
	if a != 2 || b != 1 {
		t.Errorf("Variables not swapped properly. Values were a: %d, b: %d", a, b)
	}
}

func BenchmarkSwapIntegers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Should not incur any cost for swapping these two numbers
		a, b := i-1, i-2
		swapTwoIntegers(&a, &b)
	}
}
