# Publishing intervals-cli

This document covers the one-time setup required to publish releases and the
steps to cut a new release.

Releases are automated via [goreleaser](https://goreleaser.com) and GitHub
Actions. Pushing a semver tag triggers a build for macOS (amd64/arm64) and
Linux (amd64/arm64), creates a GitHub Release with attached archives, and
updates a Homebrew formula so users can install via `brew`.

## One-Time Setup

### 1. Push the repo to GitHub

```bash
gh repo create runningcode/intervals-cli --public --source=. --push
```

### 2. Create a Homebrew tap repo

A tap is a public GitHub repo named `homebrew-tap`. goreleaser will push
updated formulas to it on each release.

```bash
gh repo create runningcode/homebrew-tap --public --clone
cd homebrew-tap
mkdir Formula
echo "# Homebrew Tap" > README.md
echo "Install: \`brew tap runningcode/tap && brew install intervals-cli\`" >> README.md
git add . && git commit -m "init: create homebrew tap" && git push
cd -
```

### 3. Create a personal access token for goreleaser

goreleaser needs write access to the tap repo to push formula updates.

1. Go to **GitHub → Settings → Developer settings → Personal access tokens → Fine-grained tokens**
2. Create a new token with:
   - **Repository access:** `runningcode/homebrew-tap` only
   - **Permissions:** Contents → Read and write
3. Copy the token value

### 4. Add repository secrets

In the intervals-cli repo on GitHub, go to **Settings → Secrets and variables → Actions** and add:

| Secret | Value |
|--------|-------|
| `HOMEBREW_TAP_TOKEN` | The PAT you created in step 3 |

`GITHUB_TOKEN` is provided automatically by GitHub Actions — no action needed.

### 5. Update .goreleaser.yaml with your GitHub username

Replace `runningcode` in `.goreleaser.yaml` with your actual GitHub
username so goreleaser knows where to push the Homebrew formula.

## Releasing

Once setup is complete, releases are a single command:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions release workflow triggers automatically, runs tests, builds
binaries for all platforms, creates a GitHub Release, and updates the Homebrew
formula.

After the release completes, users can install with:

```bash
brew tap runningcode/tap
brew install intervals-cli
```

## Version Guidelines

Follow [semver](https://semver.org). Because agents may parse the JSON output
programmatically, treat output schema changes as breaking:

| Change | Version bump |
|--------|-------------|
| Bug fixes, no output changes | Patch (1.0.x) |
| New commands or flags, existing output unchanged | Minor (1.x.0) |
| Removed commands/flags, or output schema changes | Major (x.0.0) |

## Testing a Release Locally

To verify the goreleaser config is valid before tagging:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

`--snapshot` builds without a git tag and skips publishing. Artifacts land in
`./dist/`.
