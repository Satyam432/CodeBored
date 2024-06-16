package main

import (
	"fmt"
	gemini "geminiProject/gemini"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	//Initialize gemini api
	gemini.GeminiInit()

}
