package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Please use: go run cmd/remembrall/main.go")
	fmt.Println("Or build with: go build -o remembrall cmd/remembrall/main.go")
	os.Exit(1)
}