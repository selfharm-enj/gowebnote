// Встраивание html/* в embed.FS.
package templates

import "embed"

// HTML содержит все файлы из каталога 'html/' попадающие под шаблон '*.tmpl.html'.
//
//go:embed html/*.tmpl.html
var HtmlFS embed.FS
