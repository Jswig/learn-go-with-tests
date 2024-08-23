package main

import "fmt"

const (
    englishHelloPrefix = "Hello, "
    spanishHelloPrefix = "Hola, "
    frenchHelloPrefix = "Bonjour, "
)

func Hello(name string, language string) string { 
    if name == "" {
		name = "world"
	}
    prefix := greetingPrefix(language)
    return prefix + name
}

func greetingPrefix(language string) (prefix string) {
    switch language {
    case "French":
        prefix = frenchHelloPrefix 
    case "Spanish":
        prefix = spanishHelloPrefix
    case "English":
        prefix = englishHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return
}

func main() {
	fmt.Println(Hello("world", ""))
}
