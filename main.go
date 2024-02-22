package main

import (
	"github.com/leventogut/xec/pkg/cmd"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env and ignore errors
	_ = godotenv.Load()
	cmd.Execute()
}
