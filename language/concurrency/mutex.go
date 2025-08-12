package concurrency

import (
	"fmt"
	"sync"
)

/*
Mutex allows only one goroutine to have access to a resource at a given time.
This prevents a race condition when two or more goroutines are simultaneously reading and/or writing to the same variable.
*/
var (
	counter int
	mu      sync.Mutex
)

func increment(wg *sync.WaitGroup) {
	defer wg.Done()

	mu.Lock()   // catch mutex
	counter++   // critical section
	mu.Unlock() // releasing mutex
}

func mutexFunc() {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

/*
Мьютексы и синхронизация в Go:
1. sync.Mutex
   - Простой мьютекс для взаимного исключения.
   - Только один поток (горутина) может держать блокировку.
   - Не поддерживает повторный захват (не рекурсивный).
2. sync.RWMutex
   - Мьютекс с поддержкой разделения на чтение и запись.
   - Несколько горутин могут одновременно читать (RLock).
   - Только одна горутина может писать (Lock), блокируя чтение и другие записи.
3. sync.Once
   - Гарантирует, что функция будет выполнена только один раз (например, инициализация).
   - Безопасен для использования из нескольких горутин.
4. sync.WaitGroup
   - Не мьютекс, но используется для ожидания завершения группы горутин.
   - Часто применяется вместе с мьютексами.
5. sync.Cond
   - Условная переменная для сложной координации между горутинами.
   - Основана на мьютексе и позволяет горутинам ожидать сигнала (`Wait`, `Signal`, `Broadcast`).
6. sync/atomic
   - Пакет для низкоуровневых атомарных операций (например, атомарное инкрементирование int32/int64).
   - Альтернатива мьютексам, когда нужна высокая производительность.
*/
