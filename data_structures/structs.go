package data_structures

import "fmt"

type Address struct { // Structure can be anonymous
	Street, City, State string
}

type Person struct {
	Name    string
	Age     int
	Address // functions can contain other functions
}

func (p *Person) HaveBirthday() { // Function can belong to a structure
	p.Age++
}

func structs() {
	p := Person{
		Name: "Alice",
		Age:  30,
	}
	// or p := Person{"Alice", 30} or p := &Person{"Alice", 30}

	p1 := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			Street: "123 Elm St",
			City:   "Springfield",
			State:  "IL",
		},
	}
	fmt.Println(p.Name)
	fmt.Println(p.City)
	p.HaveBirthday()
	fmt.Println(p.Age)
	fmt.Println(p1.Name)
	fmt.Println(p1.City)
}
