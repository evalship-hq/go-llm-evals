// evalworker runs the support bot for tinyeval. In evals/support.sh:
//
//	worker go run ./cmd/evalworker
//
// The client is built once per worker, not once per test case.
package main

import (
	"context"
	"os"

	tinyeval "github.com/evalship-hq/tinyeval/go"

	"github.com/evalship-hq/go-llm-evals/support"
)

func main() {
	bot := support.New()
	// The helper counts API errors that aren't the app's fault (401, 429, 5xx, timeouts, the network)
	// as couldn't run. Any other error or a panic is the reply, so the eval fails.
	tinyeval.Worker(func(input string) (string, error) {
		// The SDK's missing-key error is internal, so the worker checks the key itself.
		if os.Getenv("ANTHROPIC_API_KEY") == "" {
			return "", tinyeval.CouldntRun("ANTHROPIC_API_KEY isn't set")
		}
		return bot.Answer(context.Background(), input)
	})
}
