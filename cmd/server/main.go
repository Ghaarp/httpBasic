package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

const (
	url           = "localhost"
	port          = "3000"
	createAddress = "/create"
	getAddress    = "/get%d"
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
	data: make(map[int]*Note),
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

func main() {
	router := chi.NewRouter()
	router.Post(createAddress, createNoteHandler)
}
