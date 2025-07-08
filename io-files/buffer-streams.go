package io_files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
bufio.Reader is a wrapper around io.Reader that buffers the input — instead of reading byte-by-byte, it reads chunks into memory.
This makes reading faster and more efficient, especially with slow sources (e.g., files, network, stdin).
*/
func bufferedReadingFromString() {
	input := "Hello\nWorld\n"
	r := bufio.NewReader(strings.NewReader(input))

	data, _ := r.Peek(5) //Peeks ahead in the buffer without consuming the bytes.
	fmt.Println(string(data))

	line1, _ := r.ReadString('\n')
	line2, _ := r.ReadString('\n')

	fmt.Printf("1: %s2: %s", line1, line2)
}

func bufferedReadingFromFile() {
	f, _ := os.Open("example.txt")
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print("-> ", line)
	}
}

/*
| When to use                           |
| ------------------------------------- |
| You’re reading large text content     |
| The source is slow (file, stdin, net) |
| You want precise buffer control       |
| You need line-by-line or rune reading |

| Method               | Purpose                         |
| -------------------- | ------------------------------- |
| `Read(p []byte)`     | Reads into a byte slice         |
| `ReadByte()`         | Reads one byte                  |
| `ReadRune()`         | Reads a single UTF-8 rune       |
| `ReadString(delim)`  | Reads until the delimiter       |
| `ReadBytes(delim)`   | Reads bytes until the delimiter |
| `ReadLine()`         | Low-level line reading          |
| `Peek(n)`            | Peek ahead without advancing    |
| `Discard(n)`         | Skip the next `n` bytes         |
| `Buffered()`         | Number of bytes buffered        |
| `Reset(r io.Reader)` | Reuse the Reader with new input |

*/
