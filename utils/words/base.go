package words

import (
	"bufio"
	"embed"
	"strings"

	"github.com/geerew/friendle/utils"
)

//go:embed dictionary
var embedFS embed.FS

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Dictionary holds the set of five-letter words accepted for picks and guesses
type Dictionary struct {
	words map[string]struct{}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New loads the embedded word list for picker answers and guesser submissions
func New() (*Dictionary, error) {
	words, err := readWords("dictionary/words.txt")
	if err != nil {
		return nil, err
	}

	return &Dictionary{words: utils.ToSet(words)}, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValid reports whether the word is in the dictionary
func (d *Dictionary) IsValid(word string) bool {
	_, ok := d.words[strings.ToUpper(word)]

	return ok
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// readWords loads five-letter words from a newline-delimited file
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
