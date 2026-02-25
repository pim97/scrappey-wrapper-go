# Scrappey Go Wrapper

Official-style Go wrapper for the Scrappey API, aligned with the existing Python and Node.js wrappers.

## Install

```bash
go get github.com/scrappey/wrapper-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	scrappey "github.com/scrappey/wrapper-go"
)

func main() {
	client, err := scrappey.NewClient(os.Getenv("SCRAPPEY_API_KEY"), nil)
	if err != nil {
		panic(err)
	}

	res, err := client.Get(context.Background(), scrappey.RequestOptions{
		"url": "https://example.com",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.SolutionInt("statusCode"))
	fmt.Println(res.SolutionString("response"))
}
```

## API Surface

The client mirrors the same core methods from the other wrappers:

- `Request(ctx, payload)`
- `Get(ctx, options)`
- `Post(ctx, options)`
- `Put(ctx, options)`
- `Delete(ctx, options)`
- `Patch(ctx, options)`
- `CreateSession(ctx, options)`
- `DestroySession(ctx, sessionID)`
- `ListSessions(ctx, userID...)`
- `IsSessionActive(ctx, sessionID)`
- `CreateWebSocket(ctx, options)`

All option payloads use `map[string]any` to stay compatible with Scrappey's full and evolving request schema.

## Two Examples

1. Basic request: `go run ./examples/basic`
2. Session lifecycle: `go run ./examples/session`

Set your key first:

```bash
export SCRAPPEY_API_KEY="your_key_here"
```

PowerShell:

```powershell
$env:SCRAPPEY_API_KEY="your_key_here"
```

## Defaults

- Base URL: `https://publisher.scrappey.com/api/v1`
- Timeout: `5m`

## Deploying To GitHub Securely

- Never commit real keys to the repository.
- Keep keys only in environment variables or GitHub Actions secrets.
- `.env` files are ignored by git; use `.env.example` as a template.
- Secret scanning runs automatically in CI (`.github/workflows/secret-scan.yml`).

GitHub Actions example:

```yaml
env:
  SCRAPPEY_API_KEY: ${{ secrets.SCRAPPEY_API_KEY }}
```
