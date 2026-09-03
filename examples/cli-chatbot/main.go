package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/abhinayjangde/sarvam"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	client := sarvam.NewClient(os.Getenv("SARVAM_API_KEY"))
	var messages []sarvam.Message
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}

		if strings.EqualFold(query, "exit") {
			break
		}

		messages = append(messages, sarvam.Message{
			Role:    "user",
			Content: query,
		})
		resp, err := client.Chat.Completions.Create(
			context.Background(),
			sarvam.ChatCompletionRequest{
				Model:    sarvam.ModelSarvam105BConversations,
				Messages: messages,
			},
		)

		if err != nil {
			fmt.Println("Error:", err)
			messages = messages[:len(messages)-1]
			continue
		}

		if len(resp.Choices) == 0 {
			fmt.Println("No response generated")
			continue
		}

		answer := resp.Choices[0].Message.Content
		messages = append(messages, sarvam.Message{
			Role:    "assistant",
			Content: answer,
		})

		fmt.Println("> " + answer)

		if err := scanner.Err(); err != nil {
			fmt.Println("Input error:", err)
		}
	}

}
