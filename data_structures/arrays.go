package data_structures

import "fmt"

// Массивы это структуры фиксированного размера
// GO гарантирует сохранение порядка в массиве
type Person struct {
	Name string
	Age  int
}

func Arrays() {
	arr := [2]Person{ //как структура так и любой другой тип
		{"Alice", 30},
		{"Bob", 25},
	}
	fmt.Println(arr)
	printArray(arr) // Передача копии массива

	arr2 := [3]int{1, 2, 3}
	modifyArray(&arr2) //можно передать как ссылку
	fmt.Println(arr2)  // Массив изменится, выводит [100, 2, 3]

	//Также поддерживаются многомерные массивы
	var matrix [2][3]int
	matrix[0][0] = 1
	matrix[1][2] = 5
	fmt.Println(matrix)
	//Можно также работать со срезами:
	matrix2 := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println(matrix2)

	arr2[1] = 10     // изменим второй элемент массива на 10
	fmt.Println(arr) // Выведет [1 10 3]

	slice := arr[0:2]  // Срез, который включает только элементы от индекса 0 до индекса 2
	fmt.Println(slice) // Выведет [1 2]

	slice2 := arr2[:]   // Преобразуем массив в срез
	fmt.Println(slice2) // Выведет [1 2 3]
}

func printArray(arr [2]Person) {
	fmt.Println(arr)
}

func modifyArray(arr *[3]int) {
	arr[0] = 100
}
