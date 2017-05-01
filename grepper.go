package main

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

func (g *_Grepper) Grep(files []string, keywords ...string) ([]string, error) {
	matches := make([]string, 0)
	for _, file := range files {
		matched, err := g.grepSingleFile(file, keywords...)
		if err != nil {
			return nil, err
		}
		if matched == true {
			matches = append(matches, file)
		}
	}
	return matches, nil
}

func (g *_Grepper) grepSingleFile(file string, keywords ...string) (bool, error) {
	keywordMap := mapKeywords(keywords)

	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for keyword, _ := range keywordMap {
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
