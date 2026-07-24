// Command mermaid-hook is a Claude Code MessageDisplay hook. It reads the hook
// JSON on stdin; when the final assistant message contains ```mermaid fenced
// blocks, it renders each to SVG in-process (go-mermaid, pure Go, no browser)
// and returns hookSpecificOutput.displayContent with the fences replaced by
// inline data-URI images. On anything unexpected it prints nothing, so Claude
// Code displays the original text (safe, non-blocking).
package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	mermaid "github.com/zkrebbekx/go-mermaid"
)

type hookIn struct {
	Final bool   `json:"final"`
	Delta string `json:"delta"`
}

type hookOut struct {
	HookSpecificOutput struct {
		HookEventName  string `json:"hookEventName"`
		DisplayContent string `json:"displayContent"`
	} `json:"hookSpecificOutput"`
}

// matches a ```mermaid fenced block; group 1 is the diagram source.
var fence = regexp.MustCompile("(?s)```mermaid[ \t]*\n(.*?)\n```")

// encodeSVG makes an SVG safe to embed in a markdown image data-URI. A single
// left-to-right pass (NewReplacer) avoids double-encoding the % it introduces.
func encodeSVG(svg string) string {
	svg = strings.ReplaceAll(svg, "\r", "")
	svg = strings.ReplaceAll(svg, "\n", "")
	return strings.NewReplacer("%", "%25", "#", "%23", "(", "%28", ")", "%29").Replace(svg)
}

func main() {
	var in hookIn
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		return
	}
	if !in.Final || !strings.Contains(in.Delta, "```mermaid") {
		return
	}
	delta := strings.ReplaceAll(in.Delta, "\r\n", "\n")

	changed := false
	out := fence.ReplaceAllStringFunc(delta, func(block string) string {
		m := fence.FindStringSubmatch(block)
		if m == nil {
			return block
		}
		svg, err := mermaid.Render(m[1],
			mermaid.WithTheme(mermaid.Theme("default")),
			mermaid.WithPadding(16))
		if err != nil || !strings.Contains(string(svg), "<svg") {
			return block // leave original fence on render failure
		}
		changed = true
		return "![diagram](data:image/svg+xml;utf8," + encodeSVG(string(svg)) + ")"
	})
	if !changed {
		return
	}

	var o hookOut
	o.HookSpecificOutput.HookEventName = "MessageDisplay"
	o.HookSpecificOutput.DisplayContent = out
	_ = json.NewEncoder(os.Stdout).Encode(&o)
}
