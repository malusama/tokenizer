# Repository Guidelines

## Project Structure & Module Organization
The Go module `github.com/sugarme/tokenizer` keeps core orchestration in `tokenizer.go` while specific behaviors live inside sub-packages: `normalizer/` handles text cleanup rules, `pretokenizer/` splits strings, `model/` (BPE, wordpiece, unigram) and `decoder/` rebuild tokens, and `processor/` wires end-to-end pipelines. Shared helpers (cache utilities, Unicode helpers) live in `util/` and `file-util.go`. Integration points for loading Hugging Face configs sit in `pretrained/`, and runnable samples under `example/` mirror common flows such as `example/pretrained` and `example/unigram`. Tests reside next to their sources (`*_test.go`) so every change should expect a nearby companion test.

## Build, Test, and Development Commands
- `go build ./...` validates every package compiles with the current Go toolchain (1.23+ as defined in `go.mod`).
- `go test ./...` executes the unit suite (`bpe_test.go`, `pretokenizer_test.go`, etc.) and downloads any missing fixtures via `CachedPath`.
- `GO_TOKENIZER=/tmp/cache go test ./... -run Example` isolates integration-style examples so remote assets are cached in a writable directory.
- `go test ./... -coverprofile=coverage.out` refreshes the coverage artifact already tracked in the repo; inspect it via `go tool cover -func coverage.out`.
- `go run ./example/pretrained` or any folder under `example/` reproduces end-to-end tokenization pipelines for manual verification.

## Coding Style & Naming Conventions
Format all Go files with `gofmt -w` (tabs for indentation, standard Go spacing) and organize imports using `goimports` if available. Exported types and functions use PascalCase (`AddedVocabulary`, `NewBpeModel`), while internals stay lowerCamel. Keep filenames descriptive (`pretokenizer/bytelevel.go`) and ensure each package has a top-of-file doc comment explaining its role. Avoid committing generated caches or binaries; only JSON configs and Go sources belong in-tree.

## Testing Guidelines
Table-driven tests are the norm (`encoding_test.go`, `model/wordpiece/wordpiece_test.go`). Name tests `Test<Subject>_<Behavior>` and add example tests when demonstrating output is valuable. Maintain or improve coverage for any package touched, and rerun `go test ./...` before pushing. When a test needs pretrained files, require contributors to prime the cache by running the command once or by pointing `GO_TOKENIZER` to a prepared directory so CI stays deterministic.

## Commit & Pull Request Guidelines
Commit history favors short, action-oriented messages (e.g., `update pkg name to sugarme`). Follow the `<area>: <imperative change>` pattern when possible and write English summaries for clarity. Every PR should include: purpose and high-level design notes, mention of user-facing changes (new tokenizer models, cache behavior), any manual verification or command output (`go test ./...`), and reference issues or Hugging Face resources. Request review once lint, build, and tests pass locally, and attach screenshots or logs only when behavior is best illustrated visually.

## Security & Configuration Tips
The downloader stores assets under `~/.cache/tokenizer` by default; override it with `GO_TOKENIZER=/path/to/cache` to avoid polluting personal caches or leaking credentials in shared environments. Never commit downloaded model weights or tokenizer JSON files retrieved via `CachedPath`, and prefer HTTPS Hugging Face URLs checked into documentation instead of plaintext tokens.
