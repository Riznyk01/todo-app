package main

import (
	"fmt"
	"todo-app/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
