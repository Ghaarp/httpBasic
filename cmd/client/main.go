package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
)

const (
	url           = "http://localhost:"
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

func main() {
	note, err := createNoteClient()
	if err != nil {
		log.Fatal("Unable to create note", err)
	}

	log.Printf(color.RedString("Note created\n", color.GreenString("%+v", *note)))

}
