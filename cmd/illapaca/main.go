package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from illapaca CLI!")
	
	// Parse command-line arguments
	if len(os.Args) > 1 {
		fmt.Printf("You provided: %v\n", os.Args[1:])
	}
}