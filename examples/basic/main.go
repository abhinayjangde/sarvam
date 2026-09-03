package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abhinayjangde/sarvam"
)

func main() {
	apiKey := os.Getenv("SARVAM_API_KEY")

	client := sarvam.NewClient(apiKey)

	response, err := client.Chat.Completions.Create(
		context.Background(),
		sarvam.ChatCompletionRequest{
			Model: "sarvam-105b",
			Messages: []sarvam.Message{
				{
					Role:    "user",
					Content: "Explain what a goroutine.",
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Choices[0].Message.Content)
}
