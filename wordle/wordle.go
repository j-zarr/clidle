package wordle

import (
	"errors"
	words "golang-addon/week-1/golang-wordle/words"
	"strings"
)

const (
	MaxGuesses = 6
	WordSize   = 5
)

// LetterStatus can be none, correct, present, or absent
type LetterStatus int

const (
	// None = no status, not guessed yet
	None LetterStatus = iota
	// Absent = not in the word
	Absent
	// Present = in the word, but not in the correct position
	Present
	// Correct = in the correct position
	Correct
)

type WordleState struct {
	// Word is the Word that the user is trying to guess
	Word [WordSize]byte
	// guesses holds the guesses that the user has made.
	// currGuess is the index of the available slot in guesses.
	guesses   [MaxGuesses]guess
	currGuess int
}

// NewWordleState builds a new wordleState from a string.
// Pass in the word you want the user to guess.
func NewWordleState(word string) WordleState {
	w := WordleState{}
	copy(w.Word[:], word)
	return w
}

// appendGuess adds a guess to the wordleState. It returns an error
// if the guess is invalid.
func (w *WordleState) appendGuess(g guess) error {
	// Reject guesses when the max number of guesses has been reached
	if w.currGuess >= MaxGuesses {
		return errors.New("max guesses reached")
	}

	// Reject guesses that are not the correct length
	if len(g) != WordSize {
		return errors.New("invalid guess length")
	}

	// Reject guesses that are invalid words
	if !words.IsWord(g.string()) {
		return errors.New("invalid guess word")
	}

	w.guesses[w.currGuess] = g
	w.currGuess++
	return nil
}

// string converts a guess into a string
func (g *guess) string() string {
	str := ""
	for _, l := range g {
		if 'A' <= l.char && l.char <= 'Z' {
			str += string(l.char)
		}
	}
	return str
}

// isWordGuessed returns true when the latest guess is the correct word
func (w *WordleState) isWordGuessed() bool {
	wordIsCorrect := true
	for _, l := range w.guesses[w.currGuess-1] {
		if l.status != Correct {
			wordIsCorrect = false
		}
	}
	return wordIsCorrect
}

// shouldEndGame checks if the game should end
func (w *WordleState) shouldEndGame() bool {
	// End the game if the word is correct or the max guesses have been reached
	if w.isWordGuessed() || w.currGuess >= MaxGuesses {
		return true
	} else {
		return false
	}
}

type letter struct {
	char   byte
	status LetterStatus
}

// newLetter builds a new letter from a byte
func newLetter(char byte) letter {
	return letter{char: char, status: None}
}

type guess [WordSize]letter

// newGuess builds a new guess from a string
func newGuess(guessedWord string) guess {
	guess := guess{}
	for i, c := range guessedWord {
		guess[i] = newLetter(byte(c))
	}

	return guess
}

// updateLettersWithWord updates the status of the letters in the guess based on a word
func (g *guess) updateLettersWithWord(word [WordSize]byte) {
	for i := range g {
		l := &g[i]
		if l.char == word[i] {
			l.status = Correct
		} else if strings.Contains(string(word[:]), string(l.char)) {
			l.status = Present
		} else {
			l.status = Absent
		}
	}
}
