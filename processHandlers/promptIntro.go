package processhandlers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadRequest() string {
	//Start the promptReader
	// Receive input from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter input: ")
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	userInput = strings.TrimSpace(userInput)

	return userInput
}

//Call gemini client
// clientGemini := gemini.GetGeminiCLient()
// resp, err := clientGemini.GenerateContent(ctx, genai.Text(userInput))
// if err != nil {
// 	log.Fatal(err)
// }

// if len(resp.Candidates) > 0 {
// 	firstCandidate := resp.Candidates[0]
// 	content := firstCandidate.Content
// 	parts := content.Parts[0]
// 	fmt.Println("First candidate content parts:", parts)

// }

// func fetchApproach(input string) {

// }
