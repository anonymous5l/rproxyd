package main

import (
	"fmt"
	"html/template"
	"os"
)

var (
	Template = template.New("indexTemplate")
)

func init() {
	templateContent := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Index of {{ .Title }}</title>
	<style type="text/css">
		table tr { white-space: nowrap; }
		td.perms {}
		td.file-size { text-align: right; padding-left: 1em; }
		td.display-name { padding-left: 1em; }
		td.icon { width: 16px; height: 16px; display: block; }
	</style>
</head>
<body>
	<h1>Index of {{ .Title }}</h1>
	<table>
		<tbody>
			{{ range $item := .Items }}
			<tr>
				<td class="perms">
					<code>{{ $item.Permission }}</code>
				</td>
				<td class="file-size">
					<code>{{ if $item.IsDir }}DIR{{ else }}FILE{{ end }}</code>
				</td>
				<td class="file-size">
					<code>{{ $item.Size }}</code>
				</td>
				<td class="display-name">
					<a href="{{ $item.Path }}">{{ $item.Name }}</a>
				</td>
			</tr>
			{{ end }}
		</tbody>
	</table>
</body>
</html>`

	Template = template.Must(Template.Parse(templateContent))
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

type TemplateItem struct {
	Permission string
	Size       string
	Name       string
	Path       string
	IsDir      bool
}

type TemplateEntity struct {
	Title string
	Items []TemplateItem
}

func (te *TemplateEntity) SetTitle(dir string) {
	te.Title = fmt.Sprintf("Index of %s", dir)
}

func (te *TemplateEntity) AppendItem(info os.FileInfo, path string) {

	te.Items = append(te.Items, TemplateItem{
		Permission: fmt.Sprintf("(%s)", info.Mode().Perm().String()),
		Size:       ByteCountDecimal(info.Size()),
		Name:       info.Name(),
		Path:       path,
		IsDir:      info.IsDir(),
	})
}

func (te *TemplateEntity) Sort() {
	items := []TemplateItem{}

	for _, v := range te.Items {
		if v.IsDir {
			items = append(items, v)
		}
	}

	for _, v := range te.Items {
		if !v.IsDir {
			items = append(items, v)
		}
	}

	te.Items = items
}
