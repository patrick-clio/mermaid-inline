#!/usr/bin/env bash
# OS/arch-aware launcher: exec the matching prebuilt mermaid-hook binary.
# Self-locates via $0 (no env vars needed). On an unsupported platform or a
# missing binary it exits 0 with no output, so Claude Code shows the original
# text (MessageDisplay hooks can't block; silence == pass-through).
here=$(cd "$(dirname "$0")/.." && pwd) || exit 0
os=$(uname -s 2>/dev/null | tr '[:upper:]' '[:lower:]')
ext=""
case "$os" in
  mingw*|msys*|cygwin*|windows*) os=windows; ext=".exe" ;;
  darwin|linux) ;;
  *) exit 0 ;;
esac
arch=$(uname -m 2>/dev/null)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) exit 0 ;;
esac
bin="$here/bin/mermaid-hook_${os}_${arch}${ext}"
[ -x "$bin" ] || exit 0
exec "$bin"
