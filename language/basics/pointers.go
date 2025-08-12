package basics

import "fmt"

/*
&	взять адрес (указатель)	p := &x
*	перейти по адресу	    *p = 20
*/
func change(x *int) {
	*x = 99
}
func PointersTest0() {
	num := 10
	change(&num)
	fmt.Println(num)
}

// без указателя
func addTen1(x int) {
	x += 10
}

func PointersTest1() {
	num := 5
	addTen1(num)
	fmt.Println(num) // Выведет 5 — значение не изменилось
}

// с указателем
func addTen2(x *int) {
	*x += 10
}

func PointersTest2() {
	num := 5
	addTen2(&num)
	fmt.Println(num) // Выведет 15
}

// изменение полей структуры
type User struct {
	Name string
	Age  int
}

func birthday(u *User) {
	u.Age++
}

func PointersTest3() {
	user := User{Name: "Иван", Age: 30}
	birthday(&user)
	fmt.Println(user.Age)
}

// связанные структуры
type Node struct {
	Value int
	Next  *Node
}

func PointersTest4() {
	first := &Node{Value: 1}
	second := &Node{Value: 2}
	first.Next = second //указатели позволяют строить структуры с "ссылками" на другие элементы: списки, деревья и графы.
}
