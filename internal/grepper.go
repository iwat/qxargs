package internal

import (
	"bufio"
	"bytes"
	"os"
)

type _Grepper struct {
}

func newGrepper() *_Grepper {
	return &_Grepper{}
}

func (g *_Grepper) grep(file string, keywords ...string) (bool, error) {
	keywordMap := mapKeywords(keywords)

	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for keyword := range keywordMap {
			if bytes.Contains(scanner.Bytes(), []byte(keyword)) {
				delete(keywordMap, keyword)
			}
		}
		if len(keywordMap) == 0 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return len(keywordMap) == 0, nil
}

func mapKeywords(keywords []string) map[string]bool {
	keywordMap := map[string]bool{}

	// Create a map of all unique elements.
	for _, keyword := range keywords {
		keywordMap[keyword] = false
	}

	return keywordMap
}
