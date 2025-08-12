package using_memory_and_performance

import "fmt"

/*
Stack — a LIFO (Last In, First Out) data structure.
Used for temporary storage, recursion, parsing, graph traversal, and undo operations.
*/
type Stack []int

// Push — добавляет элемент в стек
func (s *Stack) Push(value int) {
	*s = append(*s, value)
}

// Pop — извлекает верхний элемент из стека
func (s *Stack) Pop() (int, bool) {
	stack := *s
	if len(stack) == 0 {
		return 0, false // стек пуст
	}
	last := stack[len(stack)-1]
	*s = stack[:len(stack)-1] // удаляем последний элемент
	return last, true
}

// Peek — возвращает верхний элемент без удаления
func (s *Stack) Peek() (int, bool) {
	stack := *s
	if len(stack) == 0 {
		return 0, false
	}
	return stack[len(stack)-1], true
}

func StackFunc() {
	var s Stack

	s.Push(10)
	s.Push(20)
	s.Push(30)

	fmt.Println(s.Peek()) // 30

	fmt.Println(s.Pop()) // 30
	fmt.Println(s.Pop()) // 20
	fmt.Println(s.Pop()) // 10
	fmt.Println(s.Pop()) // false
}
