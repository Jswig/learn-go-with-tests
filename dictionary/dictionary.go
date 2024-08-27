package dictionary

import "errors"

type Dictionary map[string]string

var (
    ErrNotFound = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("word already exists")
)

func (dictionary Dictionary) Add(word string, definition string) error {
    _ , wordInDictionary := dictionary[word]
    if wordInDictionary {
        return ErrWordExists 
    } else {
        dictionary[word] = definition
        return nil
    }
}

func (dictionary Dictionary) Search(word string) (string, error) {
	definition, wordInDictionary := dictionary[word]
	if wordInDictionary {
		return definition, nil
	} else {
		return "", ErrNotFound
	}
}

