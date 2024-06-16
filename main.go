package main

import (
	"fmt"
	"geminiProject/gemini"
	processhandlers "geminiProject/processHandlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	//Initialize gemini api
	gemini.GeminiInit()

	response := processhandlers.ReadRequest()
	fmt.Print("Response:", response)
}
