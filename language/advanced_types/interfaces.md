Конечно! Вот простой и понятный пример из реальной жизни: **платёжная система**, где разные способы оплаты реализуют один и тот же интерфейс.

---

## 🎯 Задача

Ты хочешь принимать оплату через **карту**, **PayPal** и **крипту**, но при этом обрабатывать их одинаково.

---

## 🧩 Интерфейс

```go
type PaymentMethod interface {
    Pay(amount float64) string
}
```

Мы говорим: *"Я хочу работать с чем-то, что умеет `Pay(float64)` и возвращает `string`."*

---

## 💳 Реализация 1 — Карта

```go
type CreditCard struct {
    Number string
}

func (c CreditCard) Pay(amount float64) string {
    return fmt.Sprintf("Оплачено %.2f руб. с карты %s", amount, c.Number)
}
```

---

## 💸 Реализация 2 — PayPal

```go
type PayPal struct {
    Email string
}

func (p PayPal) Pay(amount float64) string {
    return fmt.Sprintf("Оплачено %.2f руб. через PayPal (%s)", amount, p.Email)
}
```

---

## 🪙 Реализация 3 — Криптовалюта

```go
type CryptoWallet struct {
    Address string
}

func (w CryptoWallet) Pay(amount float64) string {
    return fmt.Sprintf("Оплачено %.2f USDT с кошелька %s", amount, w.Address)
}
```

---

## 🚀 Функция, принимающая интерфейс

```go
func ProcessPayment(p PaymentMethod, amount float64) {
    result := p.Pay(amount)
    fmt.Println(result)
}
```

---

## 🧪 Использование

```go
func main() {
    card := CreditCard{Number: "**** **** **** 1234"}
    paypal := PayPal{Email: "user@example.com"}
    crypto := CryptoWallet{Address: "0xabc123"}

    ProcessPayment(card, 1000)
    ProcessPayment(paypal, 500)
    ProcessPayment(crypto, 750)
}
```

---

## 📦 Результат

```
Оплачено 1000.00 руб. с карты **** **** **** 1234  
Оплачено 500.00 руб. через PayPal (user@example.com)  
Оплачено 750.00 USDT с кошелька 0xabc123
```

---

## 🧠 Что ты получаешь:

* Одна функция `ProcessPayment`, работающая **со всеми типами**.
* Можно **добавить новый метод оплаты**, не трогая остальной код.
* Поведение обёрнуто в **интерфейс `PaymentMethod`**.

---



## 🧭 В Go **НЕ принято**:

> "Один интерфейс — один файл" или "Одна структура — один файл"
> ❌ Такого строгого правила нет и **не нужно создавать по одному файлу на тип**, как в Java.

---

## ✅ В Go принято:

### 1. **Группировать по смыслу**, а не по сущности

* Один файл — **группа связанных типов, интерфейсов и функций**.
* Например:

```go
// file: storage.go
type Storage interface {
    Save(data []byte) error
    Load(key string) ([]byte, error)
}

type FileStorage struct { ... }
func (f FileStorage) Save(...) error { ... }
func (f FileStorage) Load(...) ([]byte, error) { ... }

type MemoryStorage struct { ... }
func (m *MemoryStorage) Save(...) error { ... }
```

### 2. **Имя файла соответствует домену / смыслу / компоненту**

Например:

| Файл              | Что в нём?                           |
| ----------------- | ------------------------------------ |
| `user.go`         | структура `User`, методы, валидаторы |
| `storage.go`      | интерфейс `Storage`, реализации      |
| `http_handler.go` | обработчики HTTP                     |
| `logger.go`       | интерфейс логгера и его реализации   |

---

## 📁 Структура пакета в Go (типично):

```bash
project/
├── main.go               # Точка входа
├── user/
│   ├── user.go           # struct User + методы
│   ├── repository.go     # интерфейс UserRepository
│   ├── postgres_repo.go  # реализация интерфейса
├── payment/
│   ├── payment.go        # interface PaymentMethod + типы
│   ├── credit_card.go    # реализация CreditCard
│   ├── paypal.go         # реализация PayPal
```

---

## 🔧 Почему в Go не делают "по типу — по файлу"

* Go **компилирует всё в одном пакете** вместе — не нужно разделять по одному.
* Чем **меньше файлов**, тем проще читать и навигировать по коду.
* Идея Go: **"читабельный, простой, быстрый для понимания"**.

---

## ✅ Best practices:

| Что        | Как лучше                                      |
| ---------- | ---------------------------------------------- |
| Интерфейсы | Хранить рядом с использованием или реализацией |
| Структуры  | В одном файле с методами                       |
| Пакеты     | По смыслу (auth, storage, user, payment)       |
| Файлы      | По логике, а не по типам                       |

---

## ✨ Вывод

> В Go **один файл = один смысловой блок**, а не "один тип/интерфейс".
> Это упрощает навигацию, уменьшает количество файлов и делает код компактнее.

## !!!Набор параметров у метода должен быть один, но так как мы вызываем реализацию метода у экземпляра структуры, то внутри вызванного метода мы можем использовать переменные структуры этого экземпляра структуры от которого мы вызываем метод

Типо такого:

```go
func ProcessPayment(pm PaymentMethod, amount float64) {
fmt.Println(pm.Pay(amount))
}

func main() {
card := CreditCard{Number: "****1234"}
paypal := PayPal{Email: "user@example.com"}

    ProcessPayment(card, 1000)   // Вызовет CreditCard.Pay
    ProcessPayment(paypal, 500) // Вызовет PayPal.Pay
}
```
