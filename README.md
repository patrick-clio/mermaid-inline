# mermaid-inline

Claude Code **desktop** plugin that renders mermaid code blocks inline as
diagrams — zero tool calls, no browser, no runtime dependencies. Install it
somewhere that displays inline images, like the Claude Code desktop app.

## Install

```
/plugin marketplace add patrick-clio/mermaid-inline
/plugin install mermaid-inline@patrick-clio
```

Restart Claude Code to activate.

## How it works

A `MessageDisplay` hook renders each mermaid code block to an SVG data-URI — via
a single static Go binary that embeds
[go-mermaid](https://github.com/zkrebbekx/go-mermaid) (pure Go, no headless
browser) — and swaps it into the message with `displayContent`. The transcript
keeps the original mermaid text, so context stays clean. A `SessionStart` hook
tells the model mermaid renders inline, so it reaches for diagrams on its own.

Verified on macOS and Linux.

## Dependencies

None.

## Building

Binaries are built in CI —
[`.github/workflows/build-binaries.yml`](.github/workflows/build-binaries.yml)
cross-compiles all six targets and commits them back whenever
`plugins/mermaid-inline/src/` changes (or on a manual run from the Actions tab).
To bump go-mermaid, edit `src/go.mod` and push; CI rebuilds.
