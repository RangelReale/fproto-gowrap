package fproto_gowrap

import (
	"fmt"
	"strings"
)

type Generator struct {
	builder strings.Builder
	indent  string
}

func NewGenerator() *Generator {
	return &Generator{}
}

// In Indents the output one tab stop.
func (g *Generator) In() { g.indent += "\t" }

// Out unindents the output one tab stop.
func (g *Generator) Out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[1:]
	}
}

func (g *Generator) WriteByte(c byte) {
	g.builder.WriteByte(c)
}

func (g *Generator) WriteString(s string) {
	g.builder.WriteString(s)
}

func (g *Generator) P(str ...interface{}) {
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
