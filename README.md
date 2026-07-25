# mermaid-inline

Claude Code **desktop** plugin that renders ` ```mermaid ``` ` blocks inline as
diagrams — zero tool calls, no browser, no runtime dependencies. It needs a
surface that displays inline images (the Claude Code desktop app); anywhere else
the block just stays a normal ` ```mermaid ` code block. Verified on macOS and
Linux; Windows is enabled but untested.

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

None.

## Building

Binaries are built in CI —
[`.github/workflows/build-binaries.yml`](.github/workflows/build-binaries.yml)
cross-compiles all six targets and commits them back whenever
`plugins/mermaid-inline/src/` changes (or on a manual run from the Actions tab).
To bump go-mermaid, edit `src/go.mod` and push; CI rebuilds.
