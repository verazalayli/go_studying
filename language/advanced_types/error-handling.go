package advanced_types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

/*
В языке Go ошибка представляется через интерфейс error:

	type error interface {
	    Error() string
	}
*/
func TestErrorHandling() {
	errors.New("что-то пошло не так")

	/*
		Как работает обработка ошибок?
		В Go принято, что функция возвращает два значения:
		Результат (или несколько)
		Ошибку (если есть)
	*/
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Результат:", result)
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("деление на ноль")
	}
	return a / b, nil
}

// Создание собственных ошибок
func findUser(id int) (string, error) {
	if id != 42 {
		return "", fmt.Errorf("пользователь с id %d не найден", id)
	}
	return "Alice", nil
}

// Обертывание ошибок
type Config struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	defer file.Close() //defer говорит компилятору, что функцию нужно выполнить в самом конце кода.

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать содержимое файла: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить JSON: %w", err)
	}

	return &cfg, nil
}

func LikeMain() {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		// Проверка: файл не найден?
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Файл конфигурации не найден. Создайте его.")
			return
		}

		// Проверка: ошибка декодирования JSON?
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			fmt.Printf("Ошибка синтаксиса JSON в байте %d\n", syntaxErr.Offset)
			return
		}

		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	fmt.Printf("Конфигурация загружена: %+v\n", cfg)
}

/*
panic и recover
Go не использует panic для ошибок — только в исключительных ситуациях, например:
Невозможность продолжать выполнение
Ошибка программиста (индекс за границей)

panic("невозможно продолжить")
Можно поймать с помощью recover, но это антипаттерн для обычной логики:


defer some func() {
    if r := recover(); r != nil {
        fmt.Println("Восстановление после паники:", r)
    }
}
*/
