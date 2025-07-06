package advanced_types

import (
	"fmt"
	"strings"
)

// How string looks inside
type stringStruct struct {
	data *byte
	len  int
}

// How to change
func changeFunc() {
	s := "hello"
	runes := []rune(s) // преобразуем в срез рун
	runes[0] = 'H'     // меняем первую руну
	s = string(runes)  // преобразуем обратно в строку
	fmt.Println(s)     // "Hello"
}

func stringsFunc() {
	// Example string
	s := "  Hello, Go World! Go is great.  "

	// ------------------ Search Functions ------------------

	// Checks if substring "Go" is in the string
	fmt.Println("Contains:", strings.Contains(s, "Go")) // true

	// Checks if the string starts with "  Hello"
	fmt.Println("HasPrefix:", strings.HasPrefix(s, "  Hello")) // true

	// Checks if the string ends with ".  "
	fmt.Println("HasSuffix:", strings.HasSuffix(s, ".  ")) // true

	// Finds the index of first occurrence of "Go"
	fmt.Println("Index:", strings.Index(s, "Go")) // 9

	// Finds the index of last occurrence of "Go"
	fmt.Println("LastIndex:", strings.LastIndex(s, "Go")) // 21

	// ------------------ Replace and Modify ------------------

	// Replace first 2 occurrences of "Go" with "Golang"
	fmt.Println("Replace:", strings.Replace(s, "Go", "Golang", 2))

	// Replace all occurrences of "Go" with "Golang"
	fmt.Println("ReplaceAll:", strings.ReplaceAll(s, "Go", "Golang"))

	// Convert to uppercase
	fmt.Println("ToUpper:", strings.ToUpper(s))

	// Convert to lowercase
	fmt.Println("ToLower:", strings.ToLower(s))

	// Remove all leading and trailing whitespace
	fmt.Println("TrimSpace:", strings.TrimSpace(s))

	// Remove specific characters from both ends
	fmt.Println("Trim:", strings.Trim(s, " .!")) // trims '.', ' ', '!' from both ends

	// Remove prefix
	fmt.Println("TrimPrefix:", strings.TrimPrefix(s, "  Hello,"))

	// Remove suffix
	fmt.Println("TrimSuffix:", strings.TrimSuffix(s, ".  "))

	// ------------------ Split and Join ------------------

	// Split the string by spaces
	words := strings.Split(s, " ")
	fmt.Println("Split:", words)

	// Split into N parts (at most 4)
	fmt.Println("SplitN:", strings.SplitN(s, " ", 4))

	// Join slice of strings into a single string with separator
	joined := strings.Join(words, "|")
	fmt.Println("Join:", joined)

	// Split by any whitespace (including multiple spaces, tabs, newlines)
	fields := strings.Fields(s)
	fmt.Println("Fields:", fields)

	// ------------------ Compare and Equal ------------------

	// Lexicographical comparison: -1, 0, or 1
	fmt.Println("Compare:", strings.Compare("apple", "banana")) // -1

	// Case-insensitive equality check
	fmt.Println("EqualFold:", strings.EqualFold("GoLang", "golang")) // true

	// ------------------ Repeat ------------------

	// Repeat the string 3 times
	fmt.Println("Repeat:", strings.Repeat("Go! ", 3))

	// ------------------ Builder Example ------------------

	var b strings.Builder
	b.WriteString("Go")
	b.WriteString(" is")
	b.WriteString(" fast!")
	fmt.Println("Builder:", b.String())

	// ------------------ NewReplacer Example ------------------

	// Replace multiple characters in one pass
	r := strings.NewReplacer("a", "@", "e", "3", "o", "0")
	result := r.Replace("awesome code")
	fmt.Println("Replacer:", result) // @w3s0m3 c0d3
}
