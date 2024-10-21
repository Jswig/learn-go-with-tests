package dictionary

type Dictionary map[string]string

const (
	ErrNotFound         = DictionaryErr("could not find the word you were looking for")
	ErrWordExists       = DictionaryErr("word already exists")
	ErrWordDoesNotExist = DictionaryErr("word does not exist")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

func (dictionary Dictionary) Add(word string, definition string) error {
	_, wordInDictionary := dictionary[word]
	if wordInDictionary {
		return ErrWordExists
	} else {
		dictionary[word] = definition
		return nil
	}
}

func (dictionary Dictionary) Update(word string, definition string) error {
	_, wordInDictionary := dictionary[word]
	if wordInDictionary {
		dictionary[word] = definition
		return nil
	} else {
		return ErrWordDoesNotExist
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

func (dictionary Dictionary) Delete(word string) {
	delete(dictionary, word)
}
