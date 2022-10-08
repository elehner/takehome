package main

import (
	"sync"
	"testing"
)

func TestSwapIntegers(t *testing.T) {
	a, b := 1, 2
	swapTwoIntegers(&a, &b)
	if a != 2 || b != 1 {
		t.Errorf("Variables not swapped properly. Values were a: %d, b: %d", a, b)
	}
}

// Confirms data is consistent between two processes
func TestCacheConsistencyBetweenRoutines(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)
	cache := NewCache()

	readUpdate := make(chan int)
	testKey := "key"
	testData := "test data"

	go func() {
		defer wg.Done()
		cache.Write(testKey, testData)
		readUpdate <- 1
	}()

	go func() {
		defer wg.Done()
		<-readUpdate
		updatedValue := cache.Read(testKey)
		if testData != updatedValue.(string) {
			t.Errorf("test data found was incorrect, found %v", updatedValue)
		}
	}()

	wg.Wait()
}

func BenchmarkCache(b *testing.B) {
	var wg sync.WaitGroup
	cache := NewCache()
	testKey := "key"
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Write(testKey, i)
			readVal := cache.Read(testKey)
			readVal = readVal.(int) + 1
		}()
	}
	wg.Wait()
}

func BenchmarkSwapIntegers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Should not incur any cost for swapping these two numbers
		a, b := i-1, i-2
		swapTwoIntegers(&a, &b)
	}
}
