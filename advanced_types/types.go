package advanced_types

import "fmt"

// Простые (базовые) типы
// Числовые типы (int и uint с разными разрядами)
var a int = 10
var b int8 = -128
var c int16 = 32767
var d int32 = -123456
var e int64 = 123456789
var f uint = 42
var g uint8 = 255 // то же самое, что byte
var h uint16 = 60000
var i uint32 = 4000000000
var j uint64 = 18446744073709551615
// Тип для работы с указателями (редко используется напрямую)
var ptr uintptr
// Строки и символы
var s string = "Hello"
var r rune = '世' // rune = int32
// Числа с плавающей точкой
var f32 float32 = 3.14
var f64 float64 = 3.14159265358979
// Комплексные числа
var z64 complex64 = 1 + 2i
var z128 complex128 = 3.14 + 1.5i
// Логический тип
var ok bool = true

// Составные типы
// Массивы (фиксированная длина)
var arr [3]int = [3]int{1, 2, 3}
// Слайсы (гибкие массивы)
var sl []int = []int{1, 2, 3}
sl = append(sl, 4)
// Карты (map) — ассоциативные массивы (словари)
m := map[string]int{"a": 1, "b": 2}
val := m["a"]
delete(m, "b")
v, ok := m["a"] // проверка наличия ключа

// Структуры (аналог классов)
type User struct {
	Name string
	Age  int
}
u := User{Name: "Alice", Age: 30}

// Интерфейсы (контракты поведения)
type Greeter interface {
	Greet() string
}
type Person struct {
	Name string
}

// Метод структуры
func (p Person) Greet() string {
	return "Hi, " + p.Name
}
var g Greeter = Person{Name: "Bob"}

// Функции как значения (первоклассные граждане)
type MathFunc func(int, int) int
func add(a, b int) int {
	return a + b
}
var f MathFunc = add
result := f(2, 3) // 5

// Указатели
x := 42
p := &x
*p = 100 // x теперь 100

// Каналы (chan) — для работы с горутинами
ch := make(chan int)
go func() {
	ch <- 42
}()
v := <-ch // получим 42

// Встроенный интерфейс error
type error interface {
	Error() string
}
err := fmt.Errorf("ошибка %d", 1)

// Тип time.Time из стандартного пакета
import "time"
now := time.Now()

// Пустой интерфейс — может принять значение любого типа
var any interface{}
any = 5
any = "строка"
any = []int{1, 2, 3}


// Примеры типовой нотации Go

// []int — слайс целых чисел
// [3]string — массив из 3 строк
// map[string]int — карта строк в числа
// chan bool — канал с булевыми значениями
// func(int) string — функция, принимающая int и возвращающая string
// Собственный тип на основе string с методом
type MyString string

func (s MyString) Hello() string {
	return "Hello, " + string(s)
}
var ms MyString = "World"
fmt.Println(ms.Hello()) // Hello, World

// Alias тип — просто другое имя, методы не добавить
type MyAlias = string
var alias MyAlias = "hello"
// Нельзя сделать: func (s MyAlias) Method() {...}
