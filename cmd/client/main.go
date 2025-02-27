package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

const (
	url           = "http://localhost:"
	port          = "3000"
	createAddress = "/create"
	getAddress    = "/get/%d"
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

func createNoteClient() (*Note, error) {

	newNote := NoteInfo{
		Author: gofakeit.Name(),
		Text:   gofakeit.BeerName(),
	}

	data, err := json.Marshal(newNote)
	if err != nil {
		return &Note{}, err
	}

	resp, err := http.Post(url+port+createAddress, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return &Note{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return &Note{}, err
	}

	var created Note
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return &Note{}, err
	}

	return &created, nil
}

func getNoteClient(id int64) (*Note, error) {
	address := fmt.Sprintf(url+port+getAddress, id)

	resp, err := http.Get(address)
	if err != nil {
		return &Note{}, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return &Note{}, errors.Errorf("Note not found: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return &Note{}, errors.Errorf("failed to get note: %d", resp.StatusCode)
	}

	var res Note
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return &Note{}, err
	}
	return &res, nil
}

func main() {
	note, err := createNoteClient()
	if err != nil {
		log.Fatal("Unable to create note", err)
	}

	log.Printf(color.RedString("Note created\n", color.GreenString("%+v", *note)))

	newNote, err2 := getNoteClient(note.ID)
	if err2 != nil {
		log.Fatal("failed to get note:", err2)
	}

	log.Printf(color.RedString("Result:/n", color.GreenString("%+v", *newNote)))
}
