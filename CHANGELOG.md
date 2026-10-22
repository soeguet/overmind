# Changelog

Fork of [DarthSim/overmind](https://github.com/DarthSim/overmind), branched off at `f57a6d0` (2.5.1).

## Unreleased

### Added
- zsh completion via `overmind completion zsh`: commands, flags, flag values, process names
- process names come from the running instance, or from the Procfile via the hidden `overmind procfile-names`
- `overmind signal <signal> [process...]`, HUP is now a supported signal (upstream #215)
- `overmind status` shows ports, `--json` for scripts (upstream #214)
- `overmind start -e file[,file]` loads env files for the processes (upstream #207)
- ci runs `go test`

### Changed
- `overmind status` has a new PORT column between PROCESS and PID (upstream #214)

### Fixed
- procfile entries without space after the colon, e.g. `web:rails s` (upstream #220)
- long project names no longer overflow the tmux socket path (upstream #204)
- `.overmind.sock` is removed when start fails fatally (upstream #204)
- ruby gem no longer prints the tmux path on start (upstream #219)
