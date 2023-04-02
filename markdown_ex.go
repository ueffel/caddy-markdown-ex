package markdown_ex

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/templates"
	"github.com/ueffel/mdtohtml/tasklistitem"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	mermaid "go.abhg.dev/goldmark/mermaid"
)

var (
	bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
	md      goldmark.Markdown
)

// MarkdownEx exposes the template function "markdown_ex" which uses a custom markdown renderer that
// includes
//   - the html class "task-list-item" for item if a TaskList
//   - support for mermaidJS syntax
type MarkdownEx struct{}

// CaddyModule returns the Caddy module information.
func (MarkdownEx) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.templates.functions.markdown_ex",
		New: func() caddy.Module { return &MarkdownEx{} },
	}
}

// CustomTemplateFunctions should return the mapping from custom function names to implementations.
func (m *MarkdownEx) CustomTemplateFunctions() template.FuncMap {
	return template.FuncMap{"markdown_ex": funcMarkdown}
}

// funcMarkdown uses a custom markdown renderer that includes
//   - the html class "task-list-item" for item if a TaskList
//   - support for mermaidJS syntax
func funcMarkdown(input any) (string, error) {
	inputStr := caddy.ToString(input)

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	err := md.Convert([]byte(inputStr), buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func unmarshalCaddyfile(d *caddyfile.Dispenser, _ any) (any, error) {
	mermaidJS := ""
	if d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "MermaidJS":
				if !d.NextArg() {
					return nil, d.ArgErr()
				}
				mermaidJS = d.Val()
			default:
				return nil, d.ArgErr()
			}
		}
	}

	md = goldmark.New(
		goldmark.WithExtensions(
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
			&mermaid.Extender{MermaidJS: mermaidJS},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	)

	return struct{}{}, nil
}

// init registers the caddy module and the markdown_ex directive.
func init() {
	httpcaddyfile.RegisterGlobalOption("markdown_ex", unmarshalCaddyfile)
	caddy.RegisterModule(MarkdownEx{})
}

// Interface guards.
var (
	_ caddy.Module              = (*MarkdownEx)(nil)
	_ templates.CustomFunctions = (*MarkdownEx)(nil)
)
