# mermaid-inline

Claude Code plugin that renders ` ```mermaid ` blocks inline as diagrams — zero
tool calls, no browser, no runtime dependencies.

## Install

```
/plugin marketplace add patrick-clio/mermaid-inline
/plugin install mermaid-inline@patrick-clio
```

Restart Claude Code to activate.

## How it works

A `MessageDisplay` hook renders each ` ```mermaid ` block to an SVG data-URI —
via a single static Go binary that embeds
[go-mermaid](https://github.com/zkrebbekx/go-mermaid) (pure Go, no headless
browser) — and swaps it into the message with `displayContent`. The transcript
keeps the original mermaid text, so context stays clean. A `SessionStart` hook
tells the model mermaid renders inline, so it reaches for diagrams on its own.

## Dependencies

None. Prebuilt binaries for all six OS/arch targets are bundled; `scripts/run.sh`
selects the right one by `uname`. macOS and Linux work out of the box. Windows
binaries and `scripts/run.ps1` are bundled but untested — enable with a
PowerShell `MessageDisplay` entry if you need it.

## Rebuild (only to bump go-mermaid)

```
cd plugins/mermaid-inline/src
go get github.com/zkrebbekx/go-mermaid@latest && go mod tidy
for t in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do
  os=${t%/*}; arch=${t#*/}; ext=""; [ "$os" = windows ] && ext=.exe
  GOOS=$os GOARCH=$arch go build -trimpath -ldflags "-s -w" -o "../bin/mermaid-hook_${os}_${arch}${ext}" .
done
```
