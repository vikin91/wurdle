package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type Logger interface {
	LogError(err error)
	LogInfof(format string, v ...interface{})
}

type Dictionary struct {
	lang          string
	secretsArr    []string
	validWordsMap map[string]struct{}
}

func loadSecrets(lang string, logger Logger) ([]string, error) {
	fileName := fmt.Sprintf("dicts/%s/%s-secrets.txt", strings.ToLower(lang), strings.ToLower(lang))
	arr := make([]string, 0)
	secretsFile, err := os.Open(fileName)
	if err != nil {
		return arr, errors.Wrapf(err, "cannot open file '%s'", fileName)
	}
	defer secretsFile.Close()
	scanner := bufio.NewScanner(secretsFile)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if utf8.RuneCountInString(word) == 6 && !strings.ContainsAny(".&!,1234567890-=+_}{[]", word) {
			arr = append(arr, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return arr, err
	}
	return arr, nil
}

func loadValidWords(lang string, length int, logger Logger) (map[string]struct{}, error) {
	var validWordsFile *os.File
	var optimizedFile *os.File
	var optimizedWriter *bufio.Writer
	m := make(map[string]struct{})

	fullValidFileName := fmt.Sprintf("dicts/%s/%s-valid.txt", strings.ToLower(lang), strings.ToLower(lang))
	optimizedFileName := fmt.Sprintf("dicts/%s/%s-valid-%d.txt", strings.ToLower(lang), strings.ToLower(lang), length)

	validWordsFile, err := os.Open(fullValidFileName)
	if err != nil {
		return m, err
	}
	defer validWordsFile.Close()
	logger.LogInfof("Reading file: " + fullValidFileName)

	optimizedFile, err = os.Create(optimizedFileName)
	if err != nil {
		return m, err
	}
	defer optimizedFile.Close()
	optimizedWriter = bufio.NewWriter(optimizedFile)

	scanner := bufio.NewScanner(validWordsFile)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		word = strings.TrimSpace(word)
		if utf8.RuneCountInString(word) == 6 && !strings.ContainsAny(".&!,1234567890-=+_}{[]", word) {
			_, _ = optimizedWriter.WriteString(word + "\n")
			m[word] = struct{}{}
		}
	}
	if optimizedWriter != nil {
		optimizedWriter.Flush()
	}
	if err := scanner.Err(); err != nil {
		return m, err
	}
	return m, nil
}

func NewDictionary(lang string, length int, logger Logger) (*Dictionary, error) {
	secrets, err := loadSecrets(lang, logger)
	if err != nil {
		return &Dictionary{}, errors.Wrapf(err, "loading secrets")
	}
	validWords, err := loadValidWords(lang, length, logger)
	if err != nil {
		return &Dictionary{}, errors.Wrapf(err, "loading valid words")
	}
	logger.LogInfof(fmt.Sprintf("Loaded %d secrets and %d valid words [%s]", len(secrets), len(validWords), lang))
	d := &Dictionary{
		lang:          lang,
		secretsArr:    secrets,
		validWordsMap: validWords,
	}
	return d, nil
}

func (d *Dictionary) GetSecret() string {
	return d.secretsArr[rand.Intn(len(d.secretsArr))]
}

func (d *Dictionary) IsWordValid(word string) bool {
	_, ok := d.validWordsMap[word]
	return ok
}
