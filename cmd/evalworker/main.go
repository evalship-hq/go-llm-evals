// evalworker runs the support bot for tinyeval. In evals/support.sh:
//
//	worker go run ./cmd/evalworker
//
// It reads each input up to a NUL byte and prints the reply and a NUL byte,
// so the client is built once per worker, not once per test case.
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"

	"github.com/evalship-hq/go-llm-evals/support"
)

func main() {
	bot := support.New()
	in := bufio.NewReader(os.Stdin)
	for {
		input, err := in.ReadString(0)
		if err != nil {
			return
		}
		reply, err := bot.Answer(context.Background(), strings.TrimSuffix(input, "\x00"))
		if err != nil {
			reply = failure(err)
		}
		fmt.Print(reply, "\x00")
	}
}

// failure turns an error into a reply. When the cause isn't the app (no key, a rate
// limit, a server error, a timeout or the network), tinyeval counts it as couldn't run.
// Any other error, like a bad model name or a bug in the app's code, is the reply, so the eval fails.
func failure(err error) string {
	var apierr *anthropic.Error
	if errors.As(err, &apierr) {
		reason := fmt.Sprintf("anthropic: %d %s", apierr.StatusCode, apierr.Type())
		if apierr.StatusCode == 401 || apierr.StatusCode == 429 || apierr.StatusCode >= 500 {
			return "#tinyeval couldnt_run " + reason + "\n"
		}
		return reason
	}
	reason, _, _ := strings.Cut(err.Error(), "\n")
	var neterr net.Error
	// The SDK's missing-credentials error is internal, so the worker checks the key itself.
	if errors.As(err, &neterr) || errors.Is(err, context.DeadlineExceeded) || os.Getenv("ANTHROPIC_API_KEY") == "" {
		return "#tinyeval couldnt_run " + reason + "\n"
	}
	return reason
}
