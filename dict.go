package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Dictionary struct {
	lang          string
	secretsArr    []string
	validWordsMap map[string]struct{}
}

func loadSecrets(lang string) []string {
	fname := "dicts/" + lang + "-secrets.txt"
	secretsFile, err := os.Open(fname)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot open file '%s'", fname))
	}
	defer secretsFile.Close()
	arr := make([]string, 0)
	scanner := bufio.NewScanner(secretsFile)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if len(word) == 6 && !strings.ContainsAny(".&!,1234567890-=+_}{[]", word) {
			arr = append(arr, word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return arr
}

func loadValidWords(lang string) map[string]struct{} {
	fname := "dicts/" + lang + "-valid.txt"
	validWordsFile, err := os.Open(fname)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot open file '%s'", fname))
	}
	defer validWordsFile.Close()
	m := make(map[string]struct{})
	scanner := bufio.NewScanner(validWordsFile)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if len(word) == 6 && !strings.ContainsAny(".&!,1234567890-=+_}{[]", word) {
			m[word] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return m
}

func NewDictionary(lang string) *Dictionary {
	d := &Dictionary{
		lang:          lang,
		secretsArr:    loadSecrets(lang),
		validWordsMap: loadValidWords(lang),
	}

	return d
}

func (d *Dictionary) GetSecret() string {
	return d.secretsArr[rand.Intn(len(d.secretsArr))]
}

func (d *Dictionary) IsWordValid(word string) bool {
	_, ok := d.validWordsMap[word]
	return ok
}
