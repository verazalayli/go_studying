package data_structures

import (
	"container/list"
	"fmt"
)

/*
l := list.New() / var l list.List — создать новый пустой список.
l.PushBack(v) — добавить элемент в конец, возвращает *Element.
l.PushFront(v) — добавить элемент в начало.
l.InsertBefore(v, mark *Element) — вставить перед mark.
l.InsertAfter(v, mark *Element) — вставить после mark.
l.Remove(e *Element) — удалить узел e.
l.MoveToFront(e *Element), l.MoveToBack(e *Element) — переместить узел к началу/концу.
l.Len() — длина списка.
l.Front(), l.Back() — получить первый/последний узел.
*/
/*
Когда используем:
Если вам важен быстрый произвольный доступ (O(1) по индексу).
Если количество элементов небольшое и нет частых вставок/удалений в середине.
Если нужна компактность и локальность данных.
*/
func LinkedList() {
	l := list.New()       // создаём новый список
	e1 := l.PushBack("A") // ["A"]
	l.PushBack("B")       // ["A", "B"]
	l.PushFront("Start")  // ["Start", "A", "B"]

	// Вставка перед e1
	l.InsertBefore("Before A", e1) // ["Start", "Before A", "A", "B"]

	// Обход списка
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

// Пример использования для реализации скользящего окна(скользящие средние, очереди и тп):
type Window struct {
	ll    *list.List
	limit int
}

func NewWindow(k int) *Window {
	return &Window{ll: list.New(), limit: k}
}

func (w *Window) Add(v int) {
	w.ll.PushBack(v)
	if w.ll.Len() > w.limit {
		w.ll.Remove(w.ll.Front())
	}
}

func (w *Window) Values() []int {
	out := make([]int, 0, w.ll.Len())
	for e := w.ll.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(int))
	}
	return out
}
