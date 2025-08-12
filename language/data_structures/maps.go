package data_structures

import "fmt"

/*
Map используется, когда вам нужно эффективно хранить и искать элементы по уникальному ключу.
map — это хеш-таблица, которая предоставляет операции с временем выполнения O(1) для вставки, удаления и поиска элементов.
*/
func Maps() {
	m := make(map[string]int) //Creating empty map, go give 8 bytes to empty map, keys have to be "hashable"

	m1 := make(map[string]int, 10) //Creating map with start capacity

	m2 := map[string]int{ //Initialisation
		"a": 1,
		"b": 2,
	}

	fmt.Println("map 1:", m)
	fmt.Println("map 2:", m1)
	fmt.Println("map 3:", m2)

	m["c"] = 3 //Adding new element, if this key already exists, value will be updated
	m["b"] = 1

	delete(m, "b")

	fmt.Println("map 1:", m)

	for key, value := range m { //Iteration on map, because go dont give incapsulated methods for it
		fmt.Println(key, value)
	}
}
