// Имплементация рендеринга и транспорта HTML шаблонов через HTTP.
package http

import (
	"net/http"

	template "webnote/internal/htmlpages"
	"webnote/utils"
)

// Handler отвечающий за рендеринг и отправку html шаблонов через HTTP.
type UIHandler struct {
	r *template.Renderer
	// localhost:8080 by default
	srvPort string
	// http OR https
	httpType string
}

// Создание нового экземпляра UIHandler.
func NewUIHandler(r *template.Renderer, srvPort, httpType string) *UIHandler {
	return &UIHandler{r: r, srvPort: srvPort, httpType: httpType}
}

func (h *UIHandler) GetSite() map[string]any {
	out := map[string]any{
		"Port":     h.srvPort,
		"HttpType": h.httpType,
	}
	return out
}

// Рендеринг домашней страницы.
func (h *UIHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	err := h.r.ExecuteTemplate(w, "home.tmpl.html", h.GetSite())
	if err != nil {
		utils.WriteJSONErrorInternalServer(w)
		return
	}
}

// Рендеринг страницы с одной запиской.
func (h *UIHandler) NotePage(w http.ResponseWriter, r *http.Request) {
	err := h.r.ExecuteTemplate(w, "note.tmpl.html", h.GetSite())
	if err != nil {
		utils.WriteJSONErrorInternalServer(w)
		return
	}
}

// Рендеринг страницы со всеми записками.
func (h *UIHandler) NotesPage(w http.ResponseWriter, r *http.Request) {
	err := h.r.ExecuteTemplate(w, "notes.tmpl.html", h.GetSite())
	if err != nil {
		utils.WriteJSONErrorInternalServer(w)
		return
	}
}

// Рендеринг страницы для входа в приложение.
func (h *UIHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.r.ExecuteTemplate(w, "login.tmpl.html", h.GetSite())
	if err != nil {
		utils.WriteJSONErrorInternalServer(w)
		return
	}
}

// Рендеринг страницы для регистрации в приложении.
func (h *UIHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	err := h.r.ExecuteTemplate(w, "register.tmpl.html", h.GetSite())
	if err != nil {
		utils.WriteJSONErrorInternalServer(w)
		return
	}
}

// Рендеринг страницы для создания записок.
func (h *UIHandler) CreateNotePage(w http.ResponseWriter, r *http.Request) {
	err := h.r.ExecuteTemplate(w, "createNote.tmpl.html", h.GetSite())
	if err != nil {
		utils.WriteJSONErrorInternalServer(w)
		return
	}
}
