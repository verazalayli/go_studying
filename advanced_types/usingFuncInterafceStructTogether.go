package advanced_types

import "fmt"

// --- Структура и интерфейс ---
type Dog struct {
	Name string
}

type Animal interface {
	Speak() string
}

func (d Dog) Speak() string {
	return d.Name + " says Woof!"
}

// --- Функция, принимающая интерфейс ---
func MakeItSpeak(a Animal) {
	fmt.Println(a.Speak())
}

// --- Дженерик-функция, работающая с любыми Animal ---
func MakeManySpeak[T Animal](animals []T) {
	for _, animal := range animals {
		fmt.Println(animal.Speak())
	}
}

// --- Пример использования ---
func usingFunc() {
	dog1 := Dog{Name: "Buddy"}
	dog2 := Dog{Name: "Charlie"}

	MakeItSpeak(dog1) // старый пример

	dogs := []Dog{dog1, dog2}
	MakeManySpeak(dogs) // Выведет: Buddy says Woof! \n Charlie says Woof!
}
