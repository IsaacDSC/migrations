package cmd

import "fmt"

func Help() {
	fmt.Println("Usage:")
	fmt.Println("  go run ./migrations up")
	fmt.Println("  go run ./migrations down")
	fmt.Println("  go run ./migrations new <filename>")
}
