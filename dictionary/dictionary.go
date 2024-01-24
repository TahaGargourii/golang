package dictionary

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	filePath string
	file     *os.File
	addCh    chan entryOperation
	removeCh chan entryOperation
	lock     *sync.Mutex
}

type entryOperation struct {
	word       string
	definition string
	resultCh   chan error
}

func New(filePath string) *Dictionary {
	d := &Dictionary{
		filePath: filePath,
		addCh:    make(chan entryOperation),
		removeCh: make(chan entryOperation),
		lock:     &sync.Mutex{}, 
	}

	go d.start()
	return d
}

func (d *Dictionary) start() {
	for {
		select {
		case addOp := <-d.addCh:
			err := d.addToDictionary(addOp.word, addOp.definition)
			if err != nil {
				log.Printf("Error adding to dictionary: %v", err)
			}
			addOp.resultCh <- err

		case removeOp := <-d.removeCh:
			err := d.removeFromDictionary(removeOp.word)
			removeOp.resultCh <- err
		}
	}
}

func (d *Dictionary) Close() error {
	if d.file != nil {
		err := d.file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dictionary) Add(word string, definition string) error {
	resultCh := make(chan error)
	d.addCh <- entryOperation{word, definition, resultCh}
	return <-resultCh
}

func (d *Dictionary) Remove(word string) error {
	// Check if word exists
	_, err := d.Get(word)
	if err != nil {
		return fmt.Errorf("word '%s' does not exist in dictionary", word)
	}

	resultCh := make(chan error)
	d.removeCh <- entryOperation{word, "", resultCh}
	return <-resultCh
}

func (d *Dictionary) Get(word string) (Entry, error) {
	file, err := os.Open(d.filePath)
	if err != nil {
		return Entry{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == word {
			return Entry{Definition: strings.TrimSpace(parts[1])}, nil
		}
	}

	return Entry{}, fmt.Errorf("word '%s' not found in the dictionary. \n", word)
}

func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	file, err := os.Open(d.filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	lines, err := readLines(file)
	if err != nil {
		return nil, nil, err
	}

	entries := make(map[string]Entry)
	var words []string

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			word := strings.TrimSpace(parts[0])
			definition := strings.TrimSpace(parts[1])
			entries[word] = Entry{Definition: definition}
			words = append(words, word)
		}
	}

	return words, entries, nil
}

func (d *Dictionary) addToDictionary(word string, definition string) error {
	if len(word) < 1 || len(word) > 50 {
		return fmt.Errorf("word must be between 1 and 50 characters")
	}
	if len(definition) < 1 || len(definition) > 200 {
		return fmt.Errorf("definition must be between 1 and 200 characters")
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	file, err := os.OpenFile(d.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	entryStr := fmt.Sprintf("%s: %s\n", word, definition)
	_, err = file.WriteString(entryStr)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dictionary) removeFromDictionary(word string) error {
	file, err := os.Open(d.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lines, err := readLines(file)
	if err != nil {
		return err
	}

	for i, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == word {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	file, err = os.Create(d.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func readLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
