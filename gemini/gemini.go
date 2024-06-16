package gemini

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var Geminimodel *genai.GenerativeModel

func GeminiInit() {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	api_key := os.Getenv("API_KEY")

	// Set up your API key
	// Create a new Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(api_key))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with most use cases
	Geminimodel = client.GenerativeModel("gemini-1.5-flash")
	if Geminimodel == nil {
		log.Fatal("Failed to load Gemini model")
	} else {
		log.Println("Gemini model loaded successfully")
	}
}

func GetGeminiCLient() *genai.GenerativeModel {
	if Geminimodel == nil {
		GeminiInit()
	}
	return Geminimodel

}
