// Имплементация html/template для работы с SSR.
package template

import (
	"html/template"
	"net/http"

	htmlfs "webnote/internal/htmlpages/templates"
)

// HTML.
type Renderer struct {
	t *template.Template
}

// Создает новый экземпляр HTMLRenderer.
// По умолчанию парсит html файлы по шаблону "html/*.tmpl.html" из htmlfs.HtmlFS.
func NewRenderer() (*Renderer, error) {
	t, err := template.ParseFS(
		htmlfs.HtmlFS,
		"html/*.tmpl.html",
	)
	if err != nil {
		return nil, err
	}
	return &Renderer{t: t}, nil
}

// Выполняет шаблон с именем 'name' в 'http.ResponseWriter' с данными `data`.
func (r *Renderer) ExecuteTemplate(w http.ResponseWriter, name string, data any) error {
	return r.t.ExecuteTemplate(w, name, data)
}
