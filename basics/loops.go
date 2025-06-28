package basics

import "fmt"

func Loops() {
	//classic for
	for k := 0; k < 4; k++ {
		fmt.Println(k)
	}

	//Range-based for
	name := "GOLANG"
	for i, s := range name {
		fmt.Printf("%d -> %c\n", i, s)
		if i == 5 {
			break //end cycle
		}
	}

	//cycle using continue
	counter := 2
	for counter < 4 {
		counter++
		if counter < 2 {
			continue //skip other lines in this cycle and start new cycle
		}
		fmt.Println(counter)
	}

	//go uses for instead of while
	count := 3
	for count < 9 {
		fmt.Println(count)
		count++
	}

	//infinite loop
	flag := 4
	for {
		flag++
		fmt.Println(flag)
	}
}
