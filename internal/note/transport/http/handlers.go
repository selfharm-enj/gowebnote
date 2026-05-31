// HTTP handlers для работы с сервисом NoteService.
package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"webnote/internal/note/usecase"
	"webnote/utils"
)

// Handler имплементирующий NoteService.
type NoteHandler struct {
	service *usecase.NoteService
}

// Создание нового экземпляра 'NoteHandler'.
func NewNoteHandler(s *usecase.NoteService) *NoteHandler {
	return &NoteHandler{service: s}
}

// Создание записки.
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	const op = "note.transtport.http.CreateNote"
	var req usecase.CreateNoteReq
	if err := utils.ReadJSON(r, &req); err != nil {
		slog.Debug("", "msg", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	if err := req.Validate(); err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	id, err := h.service.CreateNote(r.Context(), req)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	const op = "note.transtport.http.GetNote"
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	note, err := h.service.GetNote(r.Context(), int64(id))
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	res := usecase.NoteRes{
		ID:        note.ID,
		UserID:    note.UserID,
		Name:      note.Name,
		Text:      note.Text,
		CreatedAt: note.CreatedAt,
	}
	utils.WriteJSON(w, http.StatusOK, &res)
}

// Получение всех записок пользователя.
func (h *NoteHandler) GetMyNotes(w http.ResponseWriter, r *http.Request) {
	const op = "note.transtport.http.GetMyNotes"
	notes, err := h.service.GetMyNotes(r.Context())
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	utils.WriteJSON(w, http.StatusOK, &notes)
}

func (h *NoteHandler) UpdateNoteByID(w http.ResponseWriter, r *http.Request) {
	const op = "note.transtport.http.UpdateNoteByID"
	var req usecase.UpdateNoteReq
	if err := utils.ReadJSON(r, &req); err != nil {
		slog.Debug("", "msg", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	if err := req.Validate(); err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	out, err := h.service.UpdateNoteByID(r.Context(), req.ID, req)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	utils.WriteJSON(w, http.StatusOK, out)
}

func (h *NoteHandler) DeleteNoteByID(w http.ResponseWriter, r *http.Request) {
	const op = "note.transtport.http.DeleteNoteByID"
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	out, err := h.service.DeleteNoteByID(r.Context(), int64(id))
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	utils.WriteJSON(w, http.StatusOK, out)
}
