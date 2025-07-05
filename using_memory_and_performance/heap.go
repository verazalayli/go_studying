package using_memory_and_performance

/*
A heap is a memory area where objects whose lifetime cannot be predicted in advance are stored.
Unlike the stack, which is used for temporary variables inside a function and is released when exiting it,
data on the heap lives until it is collected by the garbage collector (GC).
*/
/*
Go uses "escape analysis" (analysis of going beyond the function):
If the variable "does not go beyond" the function, it is stored on the stack.
If a variable is returned from a function or captured by a closure, it ends up in the heap.
*/
func stackAllocated() int {
	x := 42 // в стеке
	return x
}

func heapAllocated() *int {
	y := 100 // в heap, т.к. возвращаем указатель
	return &y
}
