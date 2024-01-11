package dictionary

import (
	"encoding/gob"
	"errors"
	"os"
	"strings"
)

type Entry struct {
	Word       string
	Definition string
}

type Dictionary struct {
	entries map[string]Entry
	file    string
}

func New(file string) (*Dictionary, error) {
	d := &Dictionary{
		entries: make(map[string]Entry),
		file:    file,
	}
	err := d.loadFromFile()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Word: word, Definition: definition}
	d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, exists := d.entries[word]
	if !exists {
		return Entry{}, errors.New("word not found")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	var wordList []string
	for word := range d.entries {
		wordList = append(wordList, word)
	}
	return wordList, d.entries
}

func (d *Dictionary) SaveToFile() error {
	file, err := os.Create(d.file)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(d.entries)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dictionary) loadFromFile() error {
	file, err := os.Open(d.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&d.entries)
	if err != nil {
		return err
	}

	return nil
}
