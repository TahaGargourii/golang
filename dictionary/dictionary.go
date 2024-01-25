package dictionary

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	client   *redis.Client
	addCh    chan entryOperation
	removeCh chan entryOperation
	lock     *sync.Mutex
}

type entryOperation struct {
	word       string
	definition string
	resultCh   chan error
}

func New(addr string, password string, db int) *Dictionary {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Dictionary{
		client:   client,
		addCh:    make(chan entryOperation),
		removeCh: make(chan entryOperation),
		lock:     new(sync.Mutex),
	}
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

func (d *Dictionary) Add(word string, definition string) error {
	resultCh := make(chan error)
	d.addCh <- entryOperation{word, definition, resultCh}
	return <-resultCh
}

func (d *Dictionary) Remove(word string) error {
	_, err := d.Get(word)
	if err != nil {
		return fmt.Errorf("word '%s' does not exist in dictionary", word)
	}

	resultCh := make(chan error)
	d.removeCh <- entryOperation{word, "", resultCh}
	return <-resultCh
}

func (d *Dictionary) Get(word string) (Entry, error) {
	definition, err := d.client.Get(context.Background(), word).Result()
	if err == redis.Nil {
		return Entry{}, fmt.Errorf("word '%s' not found in the dictionary. \n", word)
	} else if err != nil {
		return Entry{}, err
	}

	return Entry{Definition: definition}, nil
}

func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	keys, err := d.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return nil, nil, err
	}

	entries := make(map[string]Entry)
	for _, key := range keys {
		definition, err := d.client.Get(context.Background(), key).Result()
		if err != nil {
			return nil, nil, err
		}
		entries[key] = Entry{Definition: definition}
	}

	return keys, entries, nil
}

func (d *Dictionary) addToDictionary(word string, definition string) error {
	// Validate word and definition
	if len(word) < 1 || len(word) > 50 {
		return fmt.Errorf("word must be between 1 and 50 characters")
	}
	if len(definition) < 1 || len(definition) > 200 {
		return fmt.Errorf("definition must be between 1 and 200 characters")
	}

	err := d.client.Set(context.Background(), word, definition, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (d *Dictionary) removeFromDictionary(word string) error {
	err := d.client.Del(context.Background(), word).Err()
	if err != nil {
		return err
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
