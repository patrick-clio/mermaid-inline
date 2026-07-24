# mermaid-inline

A Claude Code plugin that renders ` ```mermaid ` fenced blocks **inline as
images**, with **zero tool calls** and **no context pollution**.

## How it works

1. The model writes a normal ` ```mermaid ` code block (cheap text, low thinking).
2. A **`MessageDisplay` hook** fires as the message renders — outside the model
   loop, so there is no extra inference round-trip.
3. The hook (`mermaid-hook`, a single static Go binary that embeds
   [go-mermaid](https://github.com/zkrebbekx/go-mermaid) as a library — pure Go,
   **no headless browser**) renders each block to SVG and returns it as an
   inline data-URI via `displayContent`.
4. You see the diagram inline; the transcript keeps the original mermaid text
   (`displayContent` is display-only), so context stays clean.

No Python, no subprocess, no `jq`, no Chromium. One binary.

## Install (local marketplace)

```
/plugin marketplace add /Users/patrickodonnell/clio/mermaid-inline
/plugin install mermaid-inline@personal-tools
/reload-plugins
```

Or via settings (already wired for this machine): an `extraKnownMarketplaces`
entry of source type `directory` plus `enabledPlugins`.

## Dependencies

**None at runtime.** The plugin bundles prebuilt `mermaid-hook` binaries for all
six targets (`bin/mermaid-hook_{darwin,linux,windows}_{amd64,arm64}`); the
launcher (`scripts/run.sh`) self-selects by `uname`. Nothing to install, no Go
toolchain, no `brew`.

## Cross-OS status

| Platform | Status |
|---|---|
| macOS (arm64/amd64) | verified |
| Linux (arm64/amd64) | binaries built; same POSIX launcher path (expected to work) |
| Windows (arm64/amd64) | binaries built + `scripts/run.ps1` authored, **UNVERIFIED** — needs a Windows machine to confirm PowerShell stdin/stdout piping and hook exec |

To enable Windows, add a second `MessageDisplay` entry to `hooks/hooks.json`:

```json
{
  "hooks": [
    { "type": "command", "command": "& \"${CLAUDE_PLUGIN_ROOT}/scripts/run.ps1\"", "shell": "powershell", "timeout": 15 }
  ]
}
```

It is left out of the active config by default: on a machine without PowerShell,
a wrong-shell entry can emit a per-message hook-error notice. Add it once Windows
is verified.

## Rebuilding the binaries

Only needed to pick up a new go-mermaid version (otherwise nothing to maintain):

```
cd plugins/mermaid-inline/src
go get github.com/zkrebbekx/go-mermaid@latest && go mod tidy
for t in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do
  os=${t%/*}; arch=${t#*/}; ext=""; [ "$os" = windows ] && ext=.exe
  GOOS=$os GOARCH=$arch go build -trimpath -ldflags "-s -w" -o "../bin/mermaid-hook_${os}_${arch}${ext}" .
done
```
