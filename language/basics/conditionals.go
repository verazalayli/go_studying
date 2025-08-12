package basics

import "fmt"

func Conditionals() {
	//if-else
	age := 16
	if age < 13 {
		fmt.Println("Kid")
	} else if age < 18 {
		fmt.Println("Teenager")
	} else {
		fmt.Println("Guy")
	}

	//switch like if else
	age = 19
	var result string
	switch {
	case age < 13:
		result = "Kid"
	case age < 20:
		result = "Teenager"
	case age < 25:
		result = "Adult"
	case age < 35:
		result = "Senior"
	}
	fmt.Println("The person is a", result)

	//switch java like
	language := ""
	var devs string
	switch language {
	case "go":
		devs = "gopher"
		fallthrough //switch continue to pass(in default it breaks after matching one of the cases)
	case "rust":
		devs = "rustacean"
	case "python":
		devs = "pythonista"
	case "java":
		devs = "Duke"
	default:
		language = "javascript"
		devs = "developer"
	}
	fmt.Println("A person who codes in", language, "is called a", devs)
}
