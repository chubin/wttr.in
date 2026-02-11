package main

import (
	"fmt"

	"github.com/chubin/wttr.go/internal/generate"
)

func main() {
	err := generate.GenerateOptions()
	if err != nil {
		fmt.Println(err)
	}
}
