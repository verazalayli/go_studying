package basics

import "fmt"

// Base function
func sayHello() {
	fmt.Println("Привет!")
}

// Function with parameter
func greet(name string) {
	fmt.Println("Привет,", name)
}

// Function with returning int
func sum(a, b int) int {
	return a + b
}

// Function with returning few int(or any other types)
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// Function with named returning
func divideNames(a, b int) (quotient int, remainder int) {
	quotient = a / b
	remainder = a % b
	return // Go understand, that it have to return "quotient" and "remainder"
}

// Function belongs to structure and implement interface
type Sound struct {
	maxSound int
	minSound int
}
type Speaker interface {
	Speak() string
}

func (sound *Sound) MakeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// Anonymous function(lambda)
// filter принимает список и функцию-предикат
func filter(nums []int, predicate func(int) bool) []int {
	var result []int
	for _, n := range nums {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

// mapInts принимает список и функцию-преобразователь
func mapInts(nums []int, transformer func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = transformer(n)
	}
	return result
}

func AnonymousFunc() {
	numbers := []int{1, 2, 3, 4, 5, 6}

	// Лямбда-функция для проверки чётности
	isEven := func(n int) bool {
		return n%2 == 0
	}

	// Лямбда-функция для возведения в квадрат
	square := func(n int) int {
		return n * n
	}

	// Применение
	evenNumbers := filter(numbers, isEven)
	squaredNumbers := mapInts(evenNumbers, square)

	fmt.Println("Исходные числа:       ", numbers)
	fmt.Println("Только чётные:        ", evenNumbers)
	fmt.Println("Квадраты чётных чисел:", squaredNumbers)
}

/*
Usually, lambda functions (anonymous functions) are created inside a single function (for example, a), so that then:
immediately call them inside a;
pass it to another function (b) as an argument and use it there.
*/

// Variate parameters
func sumAnyIntegers(nums ...int) int { /* we can give any ints or slice:
	(numbers := []int{5, 10, 15}
	sum(numbers...))
	*/
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Recursive function
func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

// Recursive array traversal
func printReverse(nums []int, index int) {
	if index < 0 {
		return
	}
	fmt.Println(nums[index])
	printReverse(nums, index-1)
}

// Using defer
func deferFunc() {
	fmt.Println("Start")
	defer fmt.Println("first")
	defer fmt.Println("second")
	fmt.Println("End")
	/*
		Start
		End
		first
		second
	*/
}

// Using panic
func mayPanic() {
	panic("что-то пошло не так")
}

func main() {
	defer fmt.Println("Дефер выполнится")
	mayPanic()
	fmt.Println("Этот код не выполнится")
}

func safe() {
	defer func() {
		if r := recover(); r != nil { // Recover can be only in the defer
			fmt.Println("Паника поймана:", r)
		}
	}()
	panic("катастрофа!")
	fmt.Println("этот код не выполнится")
	// Паника поймана: катастрофа!
}
