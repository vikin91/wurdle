package ui

import "unicode"

var languageRunes map[string]map[rune]struct{}
var pangrams = map[string]string{
	"de": "Victor jagt zwölf Boxkämpfer quer über den großen Sylter Deich",
	"pl": "Stróż pchnął kość w quiz gędźb vel fax myjń",
	"en": "the quick brown fox jumps over the lazy dog",
}

var standardKeyboard = [][]rune{
	[]rune("qwertyuiop"),
	[]rune("asdfghjkl"),
	[]rune("zxcvbnm"),
}
var keyboardRows = map[string][]rune{
	"de": []rune("äüöß"),
	"pl": []rune("ęóąśłżźćń"),
	"en": []rune(""),
}

type KeyboardView struct {
	letters [][]rune
	state   map[rune]string
}

func init() {
	languageRunes = make(map[string]map[rune]struct{})
	for lang, str := range pangrams {
		languageRunes[lang] = make(map[rune]struct{})
		for _, l := range str {
			if !unicode.IsSpace(l) {
				languageRunes[lang][unicode.ToLower(l)] = struct{}{}
			}
		}

	}
}

func NewKeyboardView(lang string) *KeyboardView {
	k := &KeyboardView{
		letters: make([][]rune, 4),
	}
	k.letters[0] = keyboardRows[lang]
	k.letters[1] = standardKeyboard[0]
	k.letters[2] = standardKeyboard[1]
	k.letters[3] = standardKeyboard[2]

	return k
}
