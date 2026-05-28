package wordgame

import (
	"bufio"
	"embed"
	"strings"
	"sync"
)

//go:embed dictionary
var embedFS embed.FS

type Dictionary struct {
	AllowedGuesses map[string]struct{}
	ValidAnswers   map[string]struct{}
}

var (
	global     *Dictionary
	globalOnce sync.Once
)

func LoadDictionary() (*Dictionary, error) {
	var err error
	globalOnce.Do(func() {
		global, err = load()
	})
	return global, err
}

func load() (*Dictionary, error) {
	guesses, err := readWords("dictionary/allowed_guesses.txt")
	if err != nil {
		guesses, err = readWords("dictionary/valid_answers.txt")
		if err != nil {
			return nil, err
		}
	}
	answers, err := readWords("dictionary/valid_answers.txt")
	if err != nil {
		return nil, err
	}
	if len(guesses) == 0 {
		guesses = answers
	}
	return &Dictionary{
		AllowedGuesses: toSet(guesses),
		ValidAnswers:   toSet(answers),
	}, nil
}

func readWords(path string) ([]string, error) {
	f, err := embedFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var words []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		w := strings.ToUpper(strings.TrimSpace(sc.Text()))
		if len(w) == 5 {
			words = append(words, w)
		}
	}
	return words, sc.Err()
}

func toSet(words []string) map[string]struct{} {
	m := make(map[string]struct{}, len(words))
	for _, w := range words {
		m[w] = struct{}{}
	}
	return m
}

func (d *Dictionary) IsValidGuess(word string) bool {
	_, ok := d.AllowedGuesses[strings.ToUpper(word)]
	return ok
}

func (d *Dictionary) IsValidAnswer(word string) bool {
	_, ok := d.ValidAnswers[strings.ToUpper(word)]
	return ok
}
