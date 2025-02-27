package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
)

const (
	url           = "localhost:"
	port          = "3000"
	createAddress = "/create"
	getAddress    = "/get/{id}"
)

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteInfo struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

type SyncMap struct {
	data map[int64]*Note
	m    sync.RWMutex
}

var notes = &SyncMap{
	data: make(map[int64]*Note),
}

func createNoteHandler(writer http.ResponseWriter, request *http.Request) {
	data := &NoteInfo{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, "Failed to decode note info", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	note := &Note{
		ID:        rand.Int63(),
		Info:      *data,
		CreatedAt: time.Now(),
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(writer).Encode(note); err != nil {
		http.Error(writer, "Failed to encode note to json", http.StatusInternalServerError)
		return
	}

	notes.m.Lock()
	defer notes.m.Unlock()

	notes.data[note.ID] = note

}

func getNoteHandler(writer http.ResponseWriter, request *http.Request) {

	paramID := chi.URLParam(request, "id")
	id, err := parseNoteID(paramID)

	if err != nil {
		http.Error(writer, "Incorrect ID", http.StatusBadRequest)
		return
	}

	notes.m.RLock()
	defer notes.m.RUnlock()
	note, ok := notes.data[id]

	if !ok {
		http.Error(writer, "Note not found", http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(note); err != nil {
		http.Error(writer, "Failed to encode note to json", http.StatusInternalServerError)
		return
	}

}

func parseNoteID(str string) (int64, error) {
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func main() {

	log.Printf(color.RedString("Server started/n"))
	router := chi.NewRouter()
	router.Post(createAddress, createNoteHandler)
	router.Get(getAddress, getNoteHandler)

	if err := http.ListenAndServe(url+port, router); err != nil {
		log.Fatal("Unable to start server")
	}

}
