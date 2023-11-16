package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"estiam/dictionary"
)

func main() {
	d := dictionary.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Enter action (add, define, remove, list, exit):")
		action, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		action = strings.TrimSpace(action)

		switch action {
		case "add":
			actionAdd(d, reader)
		case "define":
			actionDefine(d, reader)
		case "remove":
			actionRemove(d, reader)
		case "list":
			actionList(d)
		case "exit":
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid action. Please enter add, define, remove, list, or exit.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Println("Enter word:")
	word, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	word = strings.TrimSpace(word)

	fmt.Println("Enter definition:")
	definition, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)
	fmt.Printf("Word '%s' added to the dictionary.\n", word)
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Println("Enter word:")
	word, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Definition of '%s': %s\n", word, entry.String())
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Println("Enter word:")
	word, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	word = strings.TrimSpace(word)

	d.Remove(word)
	fmt.Printf("Word '%s' removed from the dictionary.\n", word)
}

func actionList(d *dictionary.Dictionary) {
	words, entries := d.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word].String())
	}
}
