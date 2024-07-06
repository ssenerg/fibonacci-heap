package main

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/ssenerg/heaps/fibonacci"
)


// GenerateRandomSlice generates a slice of random integers.
func GenerateRandomSlice(size, min, max int) []int {
    slice := make([]int, size)
    for i := range slice {
        slice[i] = rand.Intn(max-min+1) + min // Generate random integer between min and max (inclusive)
    }
    return slice
}


func main() {

	for {
		slice := GenerateRandomSlice(10000000, 10, 99000)

		heap := fibonacci.NewHeap[int, int]()
		for _, key := range slice {
			heap.Insert(fibonacci.NewNode(key, key))

		}
		output := make([]int, len(slice))
		for i := 0; i < len(slice); i++ {
			minimum, err := heap.PopMin()
			if err != nil {
				fmt.Println(err)
			}
			output[i] = minimum.GetKey()
		}

		if !sort.IntsAreSorted(output) || len(output) != len(slice) {
			fmt.Println("Heap sort failed")
		} else {
			fmt.Println("Heap sort succeeded")
		}
	}

}
