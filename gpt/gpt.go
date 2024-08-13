// temp
// 後で使うので残してる
package gpt

import (
	"context"
	"fmt"
	"os"

	"github.com/otiai10/openaigo"
)

var apiKey = os.Getenv("OPENAI_API_KEY")

func NewSession() {
	client := openaigo.NewClient(apiKey)
	request := openaigo.ChatRequest{
		Model: "gpt-4o",
		Messages: []openaigo.Message{
			{Role: "user", Content: "Hello!"},
		},
	}
	ctx := context.Background()
	response, err := client.Chat(ctx, request)
	fmt.Println(response, err)
}
