package main

import "fmt"

func SayHello() string {
	return "Hello world!"
}

func main() {
	msg := SayHello()
	fmt.Println(msg)
}
