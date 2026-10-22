# overmind

Procfile process manager on top of tmux. Fork of DarthSim/overmind (stale upstream), branched off at `f57a6d0`.

## Layout
- `main.go`: cli definition (urfave/cli v1), one `setup*Cmd` per command
- `cmd_*.go`: client commands, talk to the running instance over the socket
- `start/`: the `start` command, procfile parsing, tmux, socket server (`command_center.go`)
- `completions/_overmind`: zsh completion, embedded and printed by `overmind completion zsh`

## Commands
- build: `go build`
- test: `go test ./...`
- lint: `golangci-lint run` (config in `.golangci.yml`)

## Rules
- new or renamed command, alias or flag in `main.go` → update `completions/_overmind`; `cmd_completion_test.go` fails otherwise
- every user-facing change gets a line in `CHANGELOG.md` under `Unreleased`
- keep go 1.21 compatible (`go.mod`)
- remotes: `origin` is the fork, `upstream` is DarthSim/overmind
