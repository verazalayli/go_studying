package io_files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
p []byte — slice of bytes to write.
n int — number of bytes successfully written.
err error — error if something went wrong.

type Writer interface {
	Write(p []byte) (n int, err error)
}
*/

func writingToBuffer() {
	var buf bytes.Buffer
	writer := &buf

	chunks := []string{"This ", "is ", "a ", "writer ", "example."}

	for _, chunk := range chunks {
		n, err := writer.Write([]byte(chunk))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Wrote %d bytes: %q\n", n, chunk)
	}

	fmt.Println("Full result:", buf.String())
}

/*
| Utility            | Purpose                                             |
| ------------------| --------------------------------------------------- |
| `bytes.Buffer`     | In-memory buffer that implements io.Writer          |
| `bufio.NewWriter`  | Buffered writing (flush when ready)                |
| `io.MultiWriter`   | Writes to multiple writers simultaneously           |
| `fmt.Fprintf`      | Formatted writing to any io.Writer                  |
| `io.TeeReader`     | Duplicates reads to a writer (not a writer, but related) |
*/

func writingToFile() {
	f, _ := os.Create("output.txt")
	defer f.Close()

	f.Write([]byte("Привет, файл!"))
}

func writingToString() {
	var buf bytes.Buffer
	buf.Write([]byte("Hello, "))
	buf.WriteString("world!")
	fmt.Println(buf.String())
}

func writingToConsole() {
	os.Stdout.Write([]byte("Прямо в stdout\n"))
	fmt.Fprintln(os.Stderr, "А это — в stderr")
}

func writingInParallel() {
	var buf bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &buf)

	mw.Write([]byte("Запись в консоль и в память\n"))

	fmt.Println("В памяти хранится:", buf.String())
}

func writingPlusTeeReader() {
	var buf bytes.Buffer
	reader := strings.NewReader("данные для TeeReader")
	tee := io.TeeReader(reader, &buf)

	out, _ := io.ReadAll(tee)
	fmt.Println("Прочитано:", string(out))
	fmt.Println("Скопировано в буфер:", buf.String())
}
