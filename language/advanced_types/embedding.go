package advanced_types

import "fmt"

type person struct {
	Name string
}

func (p Person) greet() {
	fmt.Println("Hello, my name is", p.Name)
}

/*
Также если определить метод greet для Employee, то этот метод переопределит метод, который принадлежит person
func (e Employee) greet() {
    fmt.Println("Hi, I'm", e.Name, "from", e.Company)
}
*/

type Employee struct {
	Person  // встроили Person
	Company string
}

// Теперь Employee наследует поведение Person
func Embedding() {
	e := Employee{
		Person:  Person{Name: "Alice"},
		Company: "Google",
	}

	e.Greet()           // Вывод: Hello, my name is Alice
	fmt.Println(e.Name) // Доступен напрямую: Alice
}

// Встраиваниевание интерфейсов
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}
