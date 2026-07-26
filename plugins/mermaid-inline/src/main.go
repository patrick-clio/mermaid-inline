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

func isFlowchart(src string) bool {
	head := strings.ToLower(strings.TrimSpace(src))
	return strings.HasPrefix(head, "flowchart") || strings.HasPrefix(head, "graph")
}

// brTag matches an HTML line-break. go-mermaid turns <br/> into real line breaks
// in flowcharts but renders it literally in other diagram types (sequence, etc.),
// leaking the raw "<br/>" text into labels. Outside flowcharts, collapse it to a
// space so the label reads as one line instead of showing the tag.
var brTag = regexp.MustCompile(`(?i)<br\s*/?>`)

func normalizeBreaks(src string) string {
	if isFlowchart(src) {
		return src
	}
	return brTag.ReplaceAllString(src, " ")
}

// edge-label forms in flowcharts: -->|label|, == label ==>, -- label -->, etc.
var (
	pipeLabel   = regexp.MustCompile(`\|\s*([^|\n]+?)\s*\|`)
	inlineLabel = regexp.MustCompile(`(?:--|==|-\.)\s+([^\n>|]+?)\s+(?:--|==|\.-)?>`)
)

// flowchartRankSep works around go-mermaid not reserving edge-label width in its
// layout (it spaces ranks by a fixed RankSep and never lengthens an edge to fit
// its label, so long labels overlap the nodes). It returns a rank gap wide enough
// for the longest edge label, or 0 when there are no labels (keep the default).
func flowchartRankSep(src string) float64 {
	longest := 0
	for _, re := range []*regexp.Regexp{pipeLabel, inlineLabel} {
		for _, m := range re.FindAllStringSubmatch(src, -1) {
			if n := len([]rune(strings.TrimSpace(m[1]))); n > longest {
				longest = n
			}
		}
	}
	if longest == 0 {
		return 0
	}
	// ~8px per char at fontSize 14, plus margin for arrowhead + node padding.
	return float64(longest)*8 + 40
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
		opts := []mermaid.Option{
			mermaid.WithTheme(mermaid.Theme("default")),
			mermaid.WithPadding(16),
		}
		if isFlowchart(m[1]) {
			if sep := flowchartRankSep(m[1]); sep > 50 {
				opts = append(opts, mermaid.WithSpacing(50, sep))
			}
		}
		svg, err := mermaid.Render(normalizeBreaks(m[1]), opts...)
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
