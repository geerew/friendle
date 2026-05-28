package wordgame

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var pepperOnce sync.Once
var pepper string

func pepperFor(dataDir string) string {
	pepperOnce.Do(func() {
		path := filepath.Join(dataDir, ".word-pepper")
		if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
			pepper = string(b)
			return
		}
		pepper = fmt.Sprintf("friendle-%d", os.Getpid())
		_ = os.WriteFile(path, []byte(pepper), 0600)
	})
	return pepper
}

func HashWord(dataDir, word string) string {
	p := pepperFor(dataDir)
	sum := sha256.Sum256([]byte(stringsToUpper(word) + ":" + p))
	return hex.EncodeToString(sum[:])
}

func stringsToUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		b[i] = c
	}
	return string(b)
}

func VerifyWord(dataDir, word, hash string) bool {
	return HashWord(dataDir, word) == hash
}
