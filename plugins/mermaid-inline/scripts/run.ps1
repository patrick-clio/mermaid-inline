# OS/arch-aware launcher (Windows). Exec the matching prebuilt binary, piping
# the hook JSON through. On failure, emit nothing so the original text is shown.
# NOTE: authored but UNVERIFIED — no Windows machine was available to test.
$ErrorActionPreference = 'SilentlyContinue'
$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
switch ($env:PROCESSOR_ARCHITECTURE) {
  'AMD64' { $a = 'amd64' }
  'ARM64' { $a = 'arm64' }
  default { exit 0 }
}
$bin = Join-Path $root ("bin/mermaid-hook_windows_" + $a + ".exe")
if (-not (Test-Path $bin)) { exit 0 }
$stdin = [Console]::In.ReadToEnd()
$stdin | & $bin
