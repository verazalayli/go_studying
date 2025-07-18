## 🧠 Почему стек предпочтительнее?

Когда данные **находятся в стеке**, это:

* **дешевле** по времени выделения/освобождения;
* **не требует участия GC**;
* автоматически освобождается по выходу из функции (нулевой overhead);
* **не вызывает фрагментации памяти**.

---

## 💥 Что значит "вытекает из функции"?

Это когда переменная внутри функции **возвращается по указателю**, или **замыкается во вложенной функции**, и **может жить дольше**, чем сама функция. Тогда Go отправит её в кучу:

```go
func stackAllocated() int {
    x := 42 // останется в стеке
    return x
}

func heapAllocated() *int {
    x := 42 // уйдёт в heap — потому что мы возвращаем *x
    return &x
}
```

---

## ❓ Тогда разве не лучше всегда копировать?

Не всегда. Действительно, **копирование значений** (например, `struct` без указателя) позволяет:

* держать данные на **стеке**;
* избежать **GC**;
* ускорить выполнение функций.

Но **всё зависит от размера структуры**:

| Размер структуры        | Копировать (значение) выгодно? | Почему                                                  |
| ----------------------- | ------------------------------ | ------------------------------------------------------- |
| Маленькая (`<=64 байт`) | ✅ Да                           | Копирование быстро, стек — дешёв                        |
| Средняя (64–512 байт)   | ⚠️ Зависит                     | Иногда лучше указатель (если структура вложена глубоко) |
| Большая (`>512 байт`)   | ❌ Нет                          | Копирование дороже, лучше работать по указателю         |

---

## 🧩 Пример: работа с `User`

```go
type User struct {
    ID    int
    Name  string
    Email string
}

// 1. Передаём по значению
func processUser(u User) {
    // вся структура скопирована
}

// 2. Передаём по указателю
func processUserPtr(u *User) {
    // структура не копируется, но может попасть в heap
}
```

👉 **Если `User` — маленький**, выгоднее `processUser(u)`
👉 **Если `User` — большой**, выгоднее `processUserPtr(&u)`

---

## ✅ Как выбрать: указатель или значение?

| Вопрос                              | Ответ                               |
| ----------------------------------- | ----------------------------------- |
| Структура маленькая и временная?    | ➤ Передавай **по значению**         |
| Структура большая (например, JSON)? | ➤ Используй **указатель**           |
| Нужно мутировать данные?            | ➤ Используй **указатель**           |
| Данные только читаются?             | ➤ Можно копировать (если не тяжело) |

---

## 📌 Вывод

* Стек = быстро, дёшево, без GC.
* Куча = медленнее, но нужна при "утечке" переменной из области видимости.
* Не **всегда** стоит избегать указателей — **иногда они дешевле, чем копирование больших структур**.
* Главное — **измерять и профилировать**.