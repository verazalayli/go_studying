package io_files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

/*
p []byte — slice into which data is read.
n int — number of bytes actually read.
err error — nil if everything went well, io.EOF if end of stream is reached, or other error.

type Reader interface {
	Read(p []byte) (n int, err error)
}
*/

func readingFiles() {
	reader := strings.NewReader("Text example")
	buf := make([]byte, 4)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read %d byte: %s\n", n, buf[:n])
	}
}

/*
| Utility           | Purpose                                          |
| ----------------- | ------------------------------------------------ |
| `bufio.NewReader` | Buffered reading (speeds up reading)             |
| `io.ReadAll`      | Reads the entire stream into memory (to the end) |
| `io.LimitReader`  | Limits the number of bytes read                  |
| `io.TeeReader`    | Duplicate the stream: reads and writes a copy    |
| `io.MultiReader`  | Reading from multiple sources in a row           |
*/

// Reading from string
func stringReading() {
	r := strings.NewReader("Hello world!")
	buf := make([]byte, 5)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Printf("Read: %s\n", buf[:n])
	}
}

// Reading from files
func fileReading() {
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 8)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Print(string(buf[:n]))
	}
}

// Reading from http response
func httpReading() {
	resp, err := http.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

// Multi reading
func multiReading() {
	r1 := strings.NewReader("Hello ")
	r2 := strings.NewReader("world!")
	r := io.MultiReader(r1, r2)

	result, _ := io.ReadAll(r)
	fmt.Println(string(result)) // Привет, мир!
}

// Limit reading
func limitReading() {
	r := strings.NewReader("1234567890")
	limited := io.LimitReader(r, 4)

	data, _ := io.ReadAll(limited)
	fmt.Println(string(data)) // 1234
}

// Reading from console
func consoleReading() {
	fmt.Print("Enter text: ")
	buf := make([]byte, 100)
	n, _ := os.Stdin.Read(buf)
	fmt.Println("You enter:", string(buf[:n]))
}
