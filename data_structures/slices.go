package data_structures

import (
	"fmt"
	"unsafe"
)

type slice struct { // Slices is a structure
	array unsafe.Pointer
	len   int
	cap   int
}

func Slices() {
	s := []int{1, 2, 3}            // Initialize slice with 3 elements
	fmt.Println(s)                 // [1 2 3]
	fmt.Println("len(s):", len(s)) // 3
	fmt.Println("cap(s):", cap(s)) // 3 or 4

	s1 := make([]int, 5) // Creating slice with 5 elements(all have zero value, because no initialization)
	fmt.Println(s1)      // [0 0 0 0 0]
	fmt.Println("len(s2):", len(s1))
	fmt.Println("cap(s2):", cap(s1))

	s3 := make([]int, 5, 10) // Len 5, cap 10
	fmt.Println(s3)          // [0 0 0 0 0]
	fmt.Println("len(s3):", len(s3))
	fmt.Println("cap(s3):", cap(s3))

	//Operations
	s = append(s, 4, 5) // Adding few elements
	fmt.Println(s)      // [1 2 3 4 5]
	fmt.Println("len(s):", len(s))
	fmt.Println("cap(s):", cap(s))

	sub := s[1:4]    // Slicing new slice from the first one with indexes from 1 to 3
	fmt.Println(sub) // [2 3 4]
	fmt.Println("len(sub):", len(sub))
	fmt.Println("cap(sub):", cap(sub))

	sub2 := s[:cap(s)] // Slicing with maximum capacity(instead of indexes)
	fmt.Println(sub2)  // [1 2 3 4]
	fmt.Println("len(sub2):", len(sub2))
	fmt.Println("cap(sub2):", cap(sub2))

	s = append(s[:2], s[3:]...) // Deleting element with index == 2
	fmt.Println(s)              // [1 2 4 5]
	s = s[:len(s)-1]            //Deleting last element
	fmt.Println(s)              // [1 2 4]
	fmt.Println("len(s):", len(s))
	fmt.Println("cap(s):", cap(s))
	fmt.Println(s[0])
}

func CalculatingCapForNewFromOldSlice() {
	s := []int{1, 2, 3, 4, 5}

	fmt.Println("cap(s):", cap(s))

	newSlice := s[3:5]

	fmt.Println("newSlice:", newSlice)
	fmt.Println("len(newSlice):", len(newSlice))
	fmt.Println("cap(newSlice):", cap(newSlice))

	s2 := make([]int, 5, 10) // Длина 5, емкость 10
	s2[0], s2[1], s2[2], s2[3], s2[4] = 1, 2, 3, 4, 5
	fmt.Println("s2:", s2)
	fmt.Println("cap(s):", cap(s2))

	newSlice2 := s2[3:5]
	fmt.Println("newSlice2:", newSlice2)
	fmt.Println("cap(newSlice2):", cap(newSlice2))

}
