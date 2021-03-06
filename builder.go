package fproto_gowrap

import (
	"fmt"
	"strings"
)

// Helper to build idented text files
type Builder struct {
	builder strings.Builder
	indent  string
}

func NewBuilder() *Builder {
	return &Builder{}
}

// In Indents the output one tab stop.
func (g *Builder) In() { g.indent += "\t" }

// Out unindents the output one tab stop.
func (g *Builder) Out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[1:]
	}
}

// Writes a single byte
func (g *Builder) WriteByte(c byte) {
	g.builder.WriteByte(c)
}

// Writes a full string
func (g *Builder) WriteString(s string) {
	g.builder.WriteString(s)
}

// Writes a list of values
func (g *Builder) P(str ...interface{}) {
	g.WriteString(g.indent)
	for _, v := range str {
		switch s := v.(type) {
		case string:
			g.WriteString(s)
		case *string:
			g.WriteString(*s)
		case bool:
			fmt.Fprintf(&g.builder, "%t", s)
		case *bool:
			fmt.Fprintf(&g.builder, "%t", *s)
		case int:
			fmt.Fprintf(&g.builder, "%d", s)
		case *int32:
			fmt.Fprintf(&g.builder, "%d", *s)
		case *int64:
			fmt.Fprintf(&g.builder, "%d", *s)
		case float64:
			fmt.Fprintf(&g.builder, "%g", s)
		case *float64:
			fmt.Fprintf(&g.builder, "%g", *s)
		default:
			panic(fmt.Sprintf("unknown type in printer: %T", v))
		}
	}
	g.WriteByte('\n')
}

// Returns the content as a string
func (g *Builder) String() string {
	return g.builder.String()
}
