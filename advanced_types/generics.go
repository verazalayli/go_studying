package advanced_types

import "fmt"

// Generic structures
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	n := len(s.items)
	item := s.items[n-1]
	s.items = s.items[:n-1]
	return item
}

// Generic func
func PrintSlice[T any](s []T) { // For any type
	for _, v := range s {
		fmt.Println(v)
	}
}

type StringOrBool interface {
	string | bool
}

func DoSomething[T StringOrBool](v T) { // For string and bool(Constraints)
	fmt.Println(v)
}

// Constraints
type Number interface {
	int | float64
}

// Sum works only with numbers
func Sum[T Number](a, b T) T {
	return a + b
}

func AdvancedTypes() {
	PrintSlice([]int{1, 2, 3})
	PrintSlice([]string{"a", "b", "c"})

	s := Stack[int]{}
	s.Push(10)
	fmt.Println(s.Pop())
}
