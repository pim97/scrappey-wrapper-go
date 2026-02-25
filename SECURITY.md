# Security Guidance

## API Keys

- Never commit real API keys.
- Use environment variables locally: `SCRAPPEY_API_KEY`.
- Use GitHub Actions secrets in CI/CD:
  - Repository `Settings` -> `Secrets and variables` -> `Actions`.
  - Create `SCRAPPEY_API_KEY`.

## GitHub Actions Usage

Reference secrets only through workflow env:

```yaml
env:
  SCRAPPEY_API_KEY: ${{ secrets.SCRAPPEY_API_KEY }}
```

## Leak Prevention

- `.env` and `.env.*` are gitignored.
- `Secret Scan` workflow runs Gitleaks on pushes and pull requests.
- Client transport errors redact query-string API keys before surfacing messages.

## If a Key Was Exposed

1. Revoke/rotate the key in Scrappey dashboard immediately.
2. Remove the key from git history (for example with `git filter-repo` or BFG).
3. Force-push cleaned history.
4. Verify with the secret scan workflow.
