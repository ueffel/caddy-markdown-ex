package markdown_ex

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/templates"
	"github.com/ueffel/mdtohtml/tasklistitem"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/mermaid"
)

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

// MarkdownEx exposes the template function "markdown_ex" which uses a custom markdown renderer that
// includes
//   - the html class "task-list-item" for item if a TaskList
//   - support for mermaidJS syntax
type MarkdownEx struct {
	md        goldmark.Markdown
	MermaidJS string `json:"mermaid_js"`
}

// CaddyModule returns the Caddy module information.
func (MarkdownEx) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.templates.functions.markdown_ex",
		New: func() caddy.Module { return &MarkdownEx{} },
	}
}

// CustomTemplateFunctions should return the mapping from custom function names to implementations.
// markdown_ex uses a custom markdown renderer that includes
//   - the html class "task-list-item" for items of a TaskList
//   - support for mermaidJS syntax
func (m *MarkdownEx) CustomTemplateFunctions() template.FuncMap {
	return template.FuncMap{"markdown_ex": func(input any) (string, error) {
		inputStr := caddy.ToString(input)

		buf := bufPool.Get().(*bytes.Buffer)
		defer bufPool.Put(buf)
		buf.Reset()

		err := m.md.Convert([]byte(inputStr), buf)
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}}
}

func (m *MarkdownEx) Provision(ctx caddy.Context) error {
	m.md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					html.WithClasses(true),
				),
			),
			extension.Linkify,
			extension.Strikethrough,
			extension.TaskList,
			tasklistitem.TaskListItemClass,
			&mermaid.Extender{MermaidURL: m.MermaidJS},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	)
	return nil
}

func (m *MarkdownEx) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "mermaid_js":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.MermaidJS = d.Val()
			default:
				return d.ArgErr()
			}
		}
	}
	return nil
}

// init registers the caddy module and the markdown_ex directive.
func init() {
	caddy.RegisterModule(MarkdownEx{})
}

// Interface guards.
var (
	_ caddy.Module              = (*MarkdownEx)(nil)
	_ templates.CustomFunctions = (*MarkdownEx)(nil)
	_ caddyfile.Unmarshaler     = (*MarkdownEx)(nil)
	_ caddy.Provisioner         = (*MarkdownEx)(nil)
)
