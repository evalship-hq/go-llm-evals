// Package support answers customer questions with Claude.
package support

import (
	"context"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"

	"github.com/evalship-hq/go-llm-evals/prompts"
)

type Bot struct {
	client anthropic.Client
}

// New builds the API client. It reads ANTHROPIC_API_KEY.
func New() *Bot {
	return &Bot{client: anthropic.NewClient()}
}

// Answer replies to one customer message.
func (b *Bot) Answer(ctx context.Context, message string) (string, error) {
	resp, err := b.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     "claude-opus-5-5",
		MaxTokens: 1024,
		// Opus 5.5 can't turn thinking off; low effort keeps a short reply quick.
		OutputConfig: anthropic.OutputConfigParam{Effort: anthropic.OutputConfigEffortLow},
		System:       []anthropic.TextBlockParam{{Text: prompts.Support}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
		},
	})
	if err != nil {
		return "", err
	}
	var reply strings.Builder
	for _, block := range resp.Content {
		if text, ok := block.AsAny().(anthropic.TextBlock); ok {
			reply.WriteString(text.Text)
		}
	}
	return reply.String(), nil
}
