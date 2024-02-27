package main

import (
	"github.com/joho/godotenv"
	"github.com/leventogut/xec/pkg/cmd"
)

func main() {
	// Load .env and ignore errors
	_ = godotenv.Load()

	cmd.Execute()
}
