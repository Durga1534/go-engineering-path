package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"markdown-note-taking/notes"
)

type HTTPHandler struct {
	service *notes.Service
}

func NewHTTPHandler(s *notes.Service) *HTTPHandler {
	return &HTTPHandler{service: s}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /notes/save", h.HandleSave)
	mux.HandleFunc("POST /notes/upload", h.HandleUpload)
	mux.HandleFunc("GET /notes", h.HandleList)
	mux.HandleFunc("GET /notes/render", h.HandleRender)
	mux.HandleFunc("POST /notes/grammar-check", h.HandleGrammarCheck)
}

func (h *HTTPHandler) HandleSave(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON body input", http.StatusBadRequest)
		return
	}
	if payload.Filename == "" || payload.Content == "" {
		http.Error(w, "Filename and content parameters are mandatory fields", http.StatusBadRequest)
		return
	}
	if err := h.service.SaveNote(payload.Filename, []byte(payload.Content)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message" : "Note saved successfully"}`))
}

func (h *HTTPHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Payload exceeds permissble upload envelope sizes", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("note")
	if err != nil {
		http.Error(w, "Missing file parameter 'note' inside multi-part transaction request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unreadable internal file structure data streams", http.StatusInternalServerError)
		return
	}
	if err := h.service.SaveNote(header.Filename, fileBytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Markdown file upload successfully"}`))
}

func (h *HTTPHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ListNotes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *HTTPHandler) HandleRender(w http.ResponseWriter, r *http.Request) {
	fileTarget := r.URL.Query().Get("file")
	if fileTarget == "" {
		http.Error(w, "Query field parameter '?file=...' must be assigned", http.StatusNotFound)
		return
	}
	htmlOutput, err := h.service.RenderToHTML(fileTarget)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			http.Error(w, "Requested record artifact not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlOutput))
}

func (h *HTTPHandler) HandleGrammarCheck(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid text payload configuartion", http.StatusBadRequest)
		return
	}
	report := h.service.CheckGrammar(payload.Text)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
