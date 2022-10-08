package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Testing cast:")
	castingInterfaceToType()
	fmt.Println("\nTesting go routines:")
	startingAndStoppingGoRoutines()
	fmt.Println("\nTesting cache:")
	simpleCachingTest()
}

func castingInterfaceToType() {
	var emptyInterface interface{} = 10

	castInterface, ok := emptyInterface.(int)
	if !ok {
		fmt.Println("Type was not applicable")
	}

	castInterface = castInterface + 5

	fmt.Printf("Properly cast interface holding the value 10 to int and added five: %d\n", castInterface)
}

func startingAndStoppingGoRoutines() {
	// Use a channel to communicate that a go process should exit
	doubleInt := make(chan int)
	printInt := make(chan int)
	quit := make(chan int)

	go func() {
		// Loop indefinitely until quit is sent,
		toBeDoubled := 1
		for {
			select {
			case <-doubleInt:
				toBeDoubled = toBeDoubled * 2
			case <-printInt:
				fmt.Printf("Int is %d\n", toBeDoubled)
			case <-quit:
				fmt.Printf("Exiting\n")
				return
			}
		}
	}()

	doubleInt <- 1
	doubleInt <- 1
	doubleInt <- 1

	printInt <- 1

	quit <- 1
}

func swapTwoIntegers(a *int, b *int) {
	// Only using this syntax to isolate this for benchmarking
	// Normally, it would just be a, b = b, a
	*a, *b = *b, *a
}

func simpleCachingTest() {
	cache := Cache{}
	cache.Init()

	// Simple RW
	cache.Write("test", "some value here")
	value := cache.Read("test").(string)
	fmt.Printf("Value was '%s'\n", value)

	// Still have to ensure it resolves concurrent reads/writes.
	// Will likely want to use a benchmark to spawn multiple processes.
}

type Cache struct {
	cache  map[string]interface{}
	rwlock sync.RWMutex
}

func (c *Cache) Init() {
	c.cache = make(map[string]interface{})
	c.rwlock = sync.RWMutex{}
}

func (c *Cache) Read(key string) interface{} {
	c.rwlock.RLock()
	cachedValue := c.cache[key]
	c.rwlock.RUnlock()
	return cachedValue
}

func (c *Cache) Write(key string, value interface{}) {
	c.rwlock.Lock()
	c.cache[key] = value
	c.rwlock.Unlock()
}
