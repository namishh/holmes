package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Question struct {
	Text            string `json:"text"`
	HashedAnswer    string `json:"hashedAnswer"`
	ImageAttachment bool   `json:"imageAttachment"`
	ImagePath       string `json:"imagePath,omitempty"`
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	questions := QuestionList{}

	for {
		var question Question

		fmt.Print("Enter a question (or press Enter to finish): ")
		scanner.Scan()
		questionText := scanner.Text()

		if questionText == "" {
			break
		}

		question.Text = questionText

		fmt.Print("Enter the answer: ")
		scanner.Scan()
		answer := scanner.Text()

		hashedAnswer, err := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error hashing answer:", err)
			continue
		}
		question.HashedAnswer = string(hashedAnswer)

		fmt.Print("Is there an image attachment? (yes/no): ")
		scanner.Scan()
		imageAttachment := strings.ToLower(scanner.Text())
		question.ImageAttachment = imageAttachment == "yes"

		if question.ImageAttachment {
			fmt.Print("Enter the path to the image: ")
			scanner.Scan()
			question.ImagePath = scanner.Text()
		}

		questions.Questions = append(questions.Questions, question)
	}

	jsonData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		fmt.Println("Error creating JSON:", err)
		return
	}

	// Write JSON to file
	err = os.WriteFile("questions.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("\nJSON output has been written to questions.json")
	fmt.Println("\nJSON content:")
	fmt.Println(string(jsonData))
}
