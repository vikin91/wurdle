package game

import (
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type Game struct {
	Dictionary *Dictionary
	secret     string

	NumGamesWon    int
	NumGamesPlayed int

	answers []string
	lang    string
}

func CreateGame(lang string, length int, logger Logger) (*Game, error) {
	g := &Game{
		Dictionary:     nil,
		NumGamesWon:    0,
		NumGamesPlayed: 0,
		answers:        make([]string, 0),
		lang:           lang,
	}
	if dict, err := NewDictionary(lang, length, logger); err != nil {
		return g, errors.Wrapf(err, "error loading dictionary '%s'", lang)
	} else {
		g.Dictionary = dict
	}
	// game should always be ready to play
	g.secret = g.Dictionary.GetSecret()
	return g, nil
}

func (g *Game) Language() string {
	return g.lang
}

func (g *Game) RevealSecret() string {
	return g.secret
}

func (g *Game) Answers() []string {
	return g.answers
}

func (g *Game) lastAnswer() string {
	if len(g.answers) == 0 {
		return ""
	}
	return g.answers[len(g.answers)-1]
}

func (g *Game) Won() bool {
	if len(g.answers) == 0 {
		return false
	}
	return g.secret == g.lastAnswer()
}

func (g *Game) CanAttempt(maxNumAttempts int) bool {
	return !g.Won() && len(g.answers) < maxNumAttempts
}

func (g *Game) Lost(maxNumAttempts int) bool {
	return !g.Won() && !g.CanAttempt(maxNumAttempts)
}

func (g *Game) NewGame() {
	g.NumGamesPlayed++
	g.secret = g.Dictionary.GetSecret()
	g.answers = make([]string, 0)
}

func (g *Game) RecordAnswer(answer string) error {
	if utf8.RuneCountInString(answer) > 6 {
		return fmt.Errorf("word '%s' too long", answer)
	}
	if utf8.RuneCountInString(answer) < 6 {
		return fmt.Errorf("word too short")
	}
	if g.Dictionary.IsWordValid(answer) {
		g.answers = append(g.answers, answer)
	} else {
		return fmt.Errorf("'%s' is not a word [%s]", answer, g.lang)
	}
	return nil
}
