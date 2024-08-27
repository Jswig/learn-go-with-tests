package dictionary

import "errors"

type Dictionary map[string]string

var ErrNotFound = errors.New("could not find the word you were looking for")

func (d Dictionary) Add(word string, definition string) {
	d[word] = definition
}

func (dictionary Dictionary) Search(word string) (string, error) {
	definition, wordInDictionary := dictionary[word]
	if wordInDictionary {
		return definition, nil
	} else {
		return "", ErrNotFound
	}
}
