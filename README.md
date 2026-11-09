# LLM evals for Go

A Go support bot whose AI call is tested on every pull request with [tinyeval](https://tinyeval.com).

## Run it

Clone this repo or use it as a template. Then:

```sh
export ANTHROPIC_API_KEY=sk-ant-...
./tinyeval
```

```
support
  ✓ refund after the window 5.8s
  ✓ missing package 5.9s
  ✓ refund inside the window 4.8s

All 3 evals passed (eval-results.json)
```

## What's here

- `support/support.go`: the bot. `Answer(ctx, message)` calls Claude with the official Go SDK.
- `prompts/support.md`: the system prompt, embedded in the binary.
- `cmd/evalworker/main.go`: builds the client once and answers each test case tinyeval sends it.
- `evals/support.sh`: three test cases, each checked by a regex, a judge or both.
- `tinyeval`: the eval runner, one bash file.
- `.github/workflows/evals.yml`: runs the evals on every pull request and push to main.

## Use it for your app

1. Replace `Answer` with your AI call. Keep the worker, so your app starts once per worker and not once per test case. Build your service's client once in the worker's `main`.
2. Write your test cases in `evals/`. The [check lines](https://github.com/evalship-hq/tinyeval#writing-evals) cover regexes, JSON, your own graders and a model judge.

When the API call fails for a reason that isn't your code, such as a missing key, a rate limit or a timeout, the worker replies `#tinyeval couldnt_run`. tinyeval then counts the test case as couldn't run, not as failed.

## In CI

Add `ANTHROPIC_API_KEY` as a repository secret. The bot and tinyeval's judge both use it. The workflow builds the code, then runs each test case 3 times. A failing eval fails the job, so it can gate merges. The results upload as the `eval-results` artifact. Without the secret, the eval step is skipped with a notice.

## Compare pull requests with main

A fresh CI run has no earlier run to compare with. [evalship](https://evalship.com) reads `eval-results.json`, compares each pull request with main, tells flaky evals from real drops and comments on the pull request.

## License

MIT
